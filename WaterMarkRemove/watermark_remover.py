import cv2
import numpy as np
import sys

def auto_detect_watermark_region(image):
    # 转换为灰度图
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
    
    # 应用自适应阈值处理
    thresh = cv2.adaptiveThreshold(
        gray, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C, cv2.THRESH_BINARY_INV, 11, 2
    )
    
    # 边缘检测
    edges = cv2.Canny(thresh, 50, 150)
    
    # 寻找轮廓
    contours, _ = cv2.findContours(edges.copy(), cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    
    # 筛选可能的水印轮廓（根据面积和宽高比）
    candidate_contours = []
    for contour in contours:
        area = cv2.contourArea(contour)
        if 1000 < area < 100000:  # 调整面积范围
            x, y, w, h = cv2.boundingRect(contour)
            aspect_ratio = w / float(h)
            if 0.5 < aspect_ratio < 10:  # 调整宽高比范围
                candidate_contours.append((x, y, w, h))
    
    # 如果找到多个候选区域，选择最可能的一个（面积最大的）
    if candidate_contours:
        return max(candidate_contours, key=lambda c: c[2]*c[3])
    
    # 如果未检测到，返回默认区域（图像右下角）
    h, w = image.shape[:2]
    return (w-200, h-50, 200, 50)

def remove_watermark(input_path, output_path, watermark_region=None):
    # 读取图像
    image = cv2.imread(input_path)
    if image is None:
        print(f"无法读取图像文件: {input_path}")
        return False

    # 智能检测水印区域
    if watermark_region is None:
        print("正在自动检测水印区域...")
        watermark_region = auto_detect_watermark_region(image)
        print(f"自动检测到水印区域: x={watermark_region[0]}, y={watermark_region[1]}, width={watermark_region[2]}, height={watermark_region[3]}")
    
    # 创建掩码（水印区域设为白色）
    mask = np.zeros(image.shape[:2], dtype=np.uint8)
    x, y, w, h = watermark_region
    mask = cv2.rectangle(mask, (x, y), (x + w, y + h), 255, -1)

    # 使用INPAINT_NS算法去除水印
    result = cv2.inpaint(image, mask, 3, cv2.INPAINT_NS)

    # 保存结果
    cv2.imwrite(output_path, result)
    print(f"水印去除完成，结果已保存至: {output_path}")
    return True

if __name__ == "__main__":
    if len(sys.argv) < 3 or len(sys.argv) > 7:
        print("用法1 (自动检测水印): python watermark_remover.py <输入图像路径> <输出图像路径>")
        print("用法2 (手动指定水印): python watermark_remover.py <输入图像路径> <输出图像路径> <x> <y> <宽度> <高度>")
        print("示例: python watermark_remover.py input.jpg output.jpg")
        sys.exit(1)

    # 解析命令行参数
    input_path = sys.argv[1]
    output_path = sys.argv[2]
    watermark_region = None

    # 如果提供了水印参数，则使用手动模式
    if len(sys.argv) == 7:
        watermark_region = (
            int(sys.argv[3]),  # x
            int(sys.argv[4]),  # y
            int(sys.argv[5]),  # width
            int(sys.argv[6])   # height
        )

    # 执行水印去除
    remove_watermark(input_path, output_path, watermark_region)