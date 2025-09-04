@echo off
REM Build script for FLV Player C++ on Windows

echo Building FLV Player C++...

REM Create build directory if it doesn't exist
if not exist "build" mkdir build

REM Change to build directory
cd build

REM Run CMake
cmake ..

REM Check if cmake succeeded
if %ERRORLEVEL% NEQ 0 (
    echo CMake configuration failed!
    cd ..
    exit /b 1
)

REM Build the project
cmake --build . --config Release

REM Check if build succeeded
if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    cd ..
    exit /b 1
)

echo Build successful!
echo Executable is located at: build\Release\flvcpp.exe
echo Run with: flvcpp.exe ^<path_to_flv_file^>