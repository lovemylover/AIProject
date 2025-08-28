#include "../include/TTSProcessor.h"
#include <iostream>
#include <fstream>
#include <cstring>

#ifdef USE_ESPEAK
// eSpeak callback function
int SynthCallback(short *wav, int numsamples, espeak_EVENT *events) {
    // This function is called when audio data is generated
    // We're not using it in this implementation, but it's required
    return 0;
}
#else
// eSpeak-ng callback function
int SynthCallback(short *wav, int numsamples, espeak_EVENT *events) {
    // This function is called when audio data is generated
    // We're not using it in this implementation, but it's required
    return 0;
}
#endif

TTSProcessor::TTSProcessor() 
    : initialized(false), voice("default"), rate(175), pitch(50), volume(100) {
}

TTSProcessor::~TTSProcessor() {
    cleanup();
}

bool TTSProcessor::initialize() {
#ifdef USE_ESPEAK
    int result = espeak_Initialize(AUDIO_OUTPUT_RETRIEVAL, 0, nullptr, espeakINITIALIZE_PHONEME_EVENTS);
    if (result == -1) {
        std::cerr << "Failed to initialize eSpeak" << std::endl;
        return false;
    }
    
    // Set callback function
    espeak_SetSynthCallback(&SynthCallback);
#else
    espeak_ng_ERROR_CONTEXT context = nullptr;
    espeak_ng_STATUS result = espeak_ng_Initialize(&context);
    if (result != ENS_OK) {
        std::cerr << "Failed to initialize eSpeak-ng" << std::endl;
        espeak_ng_ClearErrorContext(&context);
        return false;
    }
    
    // Set callback function
    espeak_SetSynthCallback(&SynthCallback);
#endif

    // Set default parameters
    setVoice(voice);
    setRate(rate);
    setPitch(pitch);
    setVolume(volume);
    
    initialized = true;
    return true;
}

void TTSProcessor::cleanup() {
    if (initialized) {
#ifdef USE_ESPEAK
        espeak_Terminate();
#else
        espeak_ng_Terminate();
#endif
        initialized = false;
    }
}

bool TTSProcessor::setText(const std::string& text) {
    if (!initialized) {
        std::cerr << "TTS processor not initialized" << std::endl;
        return false;
    }
    
#ifdef USE_ESPEAK
    int result = espeak_Synth(reinterpret_cast<const void*>(text.c_str()), 
                              text.length(), 0, POS_CHARACTER, 0, 
                              espeakCHARS_AUTO, nullptr, nullptr);
    return result == EE_OK;
#else
    espeak_ng_STATUS status = espeak_ng_Synthesize(reinterpret_cast<const void*>(text.c_str()), 
                                                   text.length(), 0, POS_CHARACTER, 0, 
                                                   espeakCHARS_AUTO, nullptr, nullptr);
    return status == ENS_OK;
#endif
}

bool TTSProcessor::synthesizeToFile(const std::string& text, const std::string& filename) {
    if (!initialized) {
        std::cerr << "TTS processor not initialized" << std::endl;
        return false;
    }
    
    // For simplicity, we'll use the synchronous approach
    // In a more advanced implementation, you might want to capture the audio data
    // and write it to a file directly
    
#ifdef USE_ESPEAK
    // Set output to file
    FILE* file = fopen(filename.c_str(), "wb");
    if (!file) {
        std::cerr << "Failed to open file for writing: " << filename << std::endl;
        return false;
    }
    
    // Close file - in a real implementation, you would capture the synthesized audio
    fclose(file);
    
    // For now, we'll just synthesize without saving to file
    int result = espeak_Synth(reinterpret_cast<const void*>(text.c_str()), 
                              text.length(), 0, POS_CHARACTER, 0, 
                              espeakCHARS_AUTO, nullptr, nullptr);
    return result == EE_OK;
#else
    // Set output to file
    FILE* file = fopen(filename.c_str(), "wb");
    if (!file) {
        std::cerr << "Failed to open file for writing: " << filename << std::endl;
        return false;
    }
    
    // Close file - in a real implementation, you would capture the synthesized audio
    fclose(file);
    
    // For now, we'll just synthesize without saving to file
    espeak_ng_STATUS status = espeak_ng_Synthesize(reinterpret_cast<const void*>(text.c_str()), 
                                                   text.length(), 0, POS_CHARACTER, 0, 
                                                   espeakCHARS_AUTO, nullptr, nullptr);
    return status == ENS_OK;
#endif
}

bool TTSProcessor::synthesizeToWav(const std::string& text, std::vector<char>& wavData) {
    if (!initialized) {
        std::cerr << "TTS processor not initialized" << std::endl;
        return false;
    }
    
    // This is a simplified implementation
    // In a full implementation, you would capture the audio samples from the callback
    // and encode them as WAV format
    
    wavData.clear();
    
#ifdef USE_ESPEAK
    int result = espeak_Synth(reinterpret_cast<const void*>(text.c_str()), 
                              text.length(), 0, POS_CHARACTER, 0, 
                              espeakCHARS_AUTO, nullptr, nullptr);
    return result == EE_OK;
#else
    espeak_ng_STATUS status = espeak_ng_Synthesize(reinterpret_cast<const void*>(text.c_str()), 
                                                   text.length(), 0, POS_CHARACTER, 0, 
                                                   espeakCHARS_AUTO, nullptr, nullptr);
    return status == ENS_OK;
#endif
}

void TTSProcessor::setVoice(const std::string& voiceName) {
    voice = voiceName;
    if (initialized) {
#ifdef USE_ESPEAK
        espeak_SetVoiceByName(voice.c_str());
#else
        espeak_SetVoiceByName(voice.c_str());
#endif
    }
}

void TTSProcessor::setRate(int wordsPerMinute) {
    rate = wordsPerMinute;
    if (initialized) {
#ifdef USE_ESPEAK
        espeak_SetParameter(espeakRATE, rate, 0);
#else
        espeak_SetParameter(espeakRATE, rate, 0);
#endif
    }
}

void TTSProcessor::setPitch(int pitchValue) {
    pitch = pitchValue;
    if (initialized) {
#ifdef USE_ESPEAK
        espeak_SetParameter(espeakPITCH, pitch, 0);
#else
        espeak_SetParameter(espeakPITCH, pitch, 0);
#endif
    }
}

void TTSProcessor::setVolume(int volumeValue) {
    volume = volumeValue;
    if (initialized) {
#ifdef USE_ESPEAK
        espeak_SetParameter(espeakVOLUME, volume, 0);
#else
        espeak_SetParameter(espeakVOLUME, volume, 0);
#endif
    }
}

std::vector<std::string> TTSProcessor::listVoices() const {
    std::vector<std::string> voices;
    
    // This is a simplified implementation
    // In a full implementation, you would query the available voices
    voices.push_back("default");
    voices.push_back("en");      // English
    voices.push_back("en-us");   // American English
    voices.push_back("en-gb");   // British English
    voices.push_back("fr");      // French
    voices.push_back("de");      // German
    voices.push_back("it");      // Italian
    voices.push_back("es");      // Spanish
    
    return voices;
}