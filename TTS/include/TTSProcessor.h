#ifndef TTS_PROCESSOR_H
#define TTS_PROCESSOR_H

#include <string>
#include <vector>

#ifdef USE_ESPEAK
#include <espeak/speak_lib.h>
#else
#include <espeak-ng/espeak_ng.h>
#endif

class TTSProcessor {
private:
    bool initialized;
    std::string voice;
    int rate;  // words per minute
    int pitch; // 0-100
    int volume; // 0-200

public:
    TTSProcessor();
    ~TTSProcessor();
    
    bool initialize();
    void cleanup();
    
    bool setText(const std::string& text);
    bool synthesizeToFile(const std::string& text, const std::string& filename);
    bool synthesizeToWav(const std::string& text, std::vector<char>& wavData);
    
    // Setters
    void setVoice(const std::string& voiceName);
    void setRate(int wordsPerMinute);
    void setPitch(int pitchValue);
    void setVolume(int volumeValue);
    
    // Getters
    std::string getVoice() const { return voice; }
    int getRate() const { return rate; }
    int getPitch() const { return pitch; }
    int getVolume() const { return volume; }
    
    // Utility functions
    std::vector<std::string> listVoices() const;
    bool isInitialized() const { return initialized; }
};

#endif // TTS_PROCESSOR_H