#include <iostream>
#include <string>
#include <vector>
#include <fstream>
#include "../include/TTSProcessor.h"

void printHelp() {
    std::cout << "Text-to-Speech (TTS) Program\n";
    std::cout << "Usage:\n";
    std::cout << "  tts [options] \"text to speak\"\n";
    std::cout << "  tts --file <filename>\n";
    std::cout << "\nOptions:\n";
    std::cout << "  --help, -h          Show this help message\n";
    std::cout << "  --rate <wpm>        Set speech rate (words per minute, default: 175)\n";
    std::cout << "  --pitch <value>     Set pitch (0-100, default: 50)\n";
    std::cout << "  --volume <value>    Set volume (0-200, default: 100)\n";
    std::cout << "  --voice <name>      Set voice (default: default)\n";
    std::cout << "  --list-voices       List available voices\n";
    std::cout << "  --file <filename>   Read text from file\n";
    std::cout << "\nExamples:\n";
    std::cout << "  tts \"Hello, world!\"\n";
    std::cout << "  tts --rate 200 --pitch 75 \"Welcome to TTS program\"\n";
    std::cout << "  tts --file input.txt\n";
    std::cout << "  tts --list-voices\n";
}

void listVoices(const TTSProcessor& tts) {
    std::cout << "Available voices:\n";
    std::vector<std::string> voices = tts.listVoices();
    for (const auto& voice : voices) {
        std::cout << "  " << voice << "\n";
    }
}

std::string readFromFile(const std::string& filename) {
    std::ifstream file(filename);
    if (!file.is_open()) {
        std::cerr << "Error: Could not open file " << filename << std::endl;
        return "";
    }
    
    std::string content((std::istreambuf_iterator<char>(file)),
                        std::istreambuf_iterator<char>());
    file.close();
    return content;
}

int main(int argc, char* argv[]) {
    if (argc < 2) {
        printHelp();
        return 1;
    }
    
    // Parse command line arguments
    std::string text;
    std::string filename;
    std::string voice = "default";
    int rate = 175;
    int pitch = 50;
    int volume = 100;
    bool listVoicesFlag = false;
    
    for (int i = 1; i < argc; i++) {
        std::string arg = argv[i];
        
        if (arg == "--help" || arg == "-h") {
            printHelp();
            return 0;
        } else if (arg == "--rate" && i + 1 < argc) {
            rate = std::stoi(argv[++i]);
        } else if (arg == "--pitch" && i + 1 < argc) {
            pitch = std::stoi(argv[++i]);
        } else if (arg == "--volume" && i + 1 < argc) {
            volume = std::stoi(argv[++i]);
        } else if (arg == "--voice" && i + 1 < argc) {
            voice = argv[++i];
        } else if (arg == "--list-voices") {
            listVoicesFlag = true;
        } else if (arg == "--file" && i + 1 < argc) {
            filename = argv[++i];
        } else if (text.empty() && arg[0] != '-') {
            text = arg;
        }
    }
    
    // Initialize TTS processor
    TTSProcessor tts;
    if (!tts.initialize()) {
        std::cerr << "Failed to initialize TTS processor" << std::endl;
        return 1;
    }
    
    // Handle list voices command
    if (listVoicesFlag) {
        listVoices(tts);
        return 0;
    }
    
    // Read text from file if specified
    if (!filename.empty()) {
        text = readFromFile(filename);
        if (text.empty()) {
            return 1;
        }
    }
    
    // Check if we have text to speak
    if (text.empty()) {
        std::cerr << "Error: No text provided\n";
        printHelp();
        return 1;
    }
    
    // Configure TTS processor
    tts.setVoice(voice);
    tts.setRate(rate);
    tts.setPitch(pitch);
    tts.setVolume(volume);
    
    // Display settings
    std::cout << "TTS Settings:\n";
    std::cout << "  Voice: " << tts.getVoice() << "\n";
    std::cout << "  Rate: " << tts.getRate() << " words per minute\n";
    std::cout << "  Pitch: " << tts.getPitch() << "\n";
    std::cout << "  Volume: " << tts.getVolume() << "\n";
    std::cout << "\nSpeaking: " << text << "\n\n";
    
    // Synthesize speech
    if (tts.setText(text)) {
        std::cout << "Speech synthesis completed successfully!\n";
    } else {
        std::cerr << "Failed to synthesize speech\n";
        return 1;
    }
    
    return 0;
}