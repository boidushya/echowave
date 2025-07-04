<h1 align="center">EchoWave</h1>

<div align="center">
  <h3>Transform audio into lyrics with AI-powered transcription</h3>
  <p>Part of the <a href="https://better-lyrics.boidu.dev" style="color: #FF0000; font-weight: bold;">better-lyrics</a> ecosystem</p>
</div>

---

## ✨ Features

🎯 **Smart Transcription** - Powered by OpenAI's Whisper AI  
📺 **YouTube Support** - Direct download and transcription from YouTube URLs  
🎵 **LRC Format** - Generates synchronized lyrics files  
🌈 **Beautiful CLI** - Colorful, animated terminal interface  
⚡ **Fast & Efficient** - Optimized for speed and accuracy  
🔧 **Configurable** - Multiple output formats and settings  
🎨 **Multi-language** - Support for 100+ languages  

## 🚀 Quick Start

**Install EchoWave:**

```bash
curl -sSL https://raw.githubusercontent.com/boidushya/echowave/main/install.sh | bash
```

*The install script automatically detects your platform, downloads dependencies, and sets up EchoWave.*

**Alternative installation methods:**
- 📦 [Download from releases](https://github.com/boidushya/echowave/releases) - Manual installation
- 🔨 Build from source: `git clone https://github.com/boidushya/echowave.git && cd echowave && ./build.sh`

### Usage

```bash
# Transcribe YouTube video
echowave https://youtube.com/watch?v=xyz

# Transcribe local audio file
echowave audio.mp3

# Custom model and language
echowave -model=medium -language=es -output=transcript audio.mp3

# Custom output directory
echowave -output-dir=transcripts https://youtube.com/watch?v=xyz

# Verbose output (show tool outputs)
echowave -verbose audio.mp3

# Disable accuracy heatmap
echowave -heatmap=false audio.mp3
```

## 🎛️ Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `-model` | Whisper model to use | `medium` |
| `-language` | Language for transcription | `en` |
| `-audio-format` | Audio format for YouTube downloads | `mp3` |
| `-output-dir` | Output directory for files | `.` |
| `-output` | Custom output filename (without extension) | Audio filename |
| `-verbose` | Show detailed output from tools | `false` |
| `-heatmap` | Show transcription accuracy heatmap | `true` |
| `-help` | Show help message | - |

### Available Whisper Models

| Model | Parameters | Speed | Accuracy |
|-------|------------|-------|----------|
| `tiny` | 39M | ⚡⚡⚡⚡⚡ | ⭐⭐ |
| `base` | 74M | ⚡⚡⚡⚡ | ⭐⭐⭐ |
| `small` | 244M | ⚡⚡⚡ | ⭐⭐⭐⭐ |
| `medium` | 749M | ⚡⚡ | ⭐⭐⭐⭐⭐ |
| `large-v3` | 1550M | ⚡ | ⭐⭐⭐⭐⭐ |

## 🌍 Supported Languages

EchoWave supports 100+ languages including:

`en` (English), `es` (Spanish), `fr` (French), `de` (German), `it` (Italian), `pt` (Portuguese), `ru` (Russian), `ja` (Japanese), `ko` (Korean), `zh` (Chinese), `ar` (Arabic), `hi` (Hindi), `th` (Thai), `vi` (Vietnamese), `id` (Indonesian), `ms` (Malay), `tl` (Filipino), `tr` (Turkish), `pl` (Polish), `nl` (Dutch), `sv` (Swedish), `da` (Danish), `no` (Norwegian), `fi` (Finnish)

[View full language list →](https://github.com/openai/whisper#available-models-and-languages)

## 📁 Output Files

EchoWave generates two files:

1. **`.json`** - Complete transcription with timestamps
2. **`.lrc`** - Synchronized lyrics file compatible with media players

### LRC Format Example
```lrc
[00:12.34] Hello world, this is a test
[00:18.56] Of the emergency broadcast system
[00:25.78] This is only a test
```

### Accuracy Heatmap

By default, EchoWave displays a color-coded visualization of transcription accuracy. Use `-heatmap=false` to disable this feature.

- 🟢 **High confidence (>0.8)** - Green text indicates words with high transcription confidence
- 🟡 **Medium confidence (0.5-0.8)** - Yellow text shows moderately confident transcription
- 🔴 **Low confidence (<0.5)** - Red text highlights uncertain or potentially incorrect words

This feature helps identify sections that may need manual review or correction.

## 🔧 Advanced Usage

### Batch Processing
```bash
# Process all MP3 files in current directory
for file in *.mp3; do
    echowave -output-dir=transcripts \"$file\"
done
```

### Custom Whisper Parameters
The tool uses optimized Whisper settings:
- `--temperature 0` for consistent output
- `--word_timestamps True` for precise timing
- `--output_format json` for structured data

### Integration with Other Tools
```bash
# Convert to SRT format using external tool
echowave song.mp3 && lrc2srt song.lrc

# Combine with video processing
echowave https://youtube.com/watch?v=xyz
ffmpeg -i video.mp4 -vf subtitles=video.lrc output.mp4
```

## 🛠️ Development

### Building from Source
```bash
# Development build
go run .

# Production build (single platform)
go build -o echowave

# Cross-compilation (all platforms)
./build.sh
```

### Testing
```bash
# Run with test audio
echowave -help

# Test YouTube download
echowave -model=tiny https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Test local file
echowave -model=base test-audio.mp3
```

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/dev`)
3. Commit your changes (`git commit -m 'feat: Added stuff'`)
4. Push to the branch (`git push origin feat/dev`)
5. Open a Pull Request

### Development Guidelines
- Write clear, documented code
- Follow Go conventions and best practices
- Add tests for new features
- Update documentation as needed
- Ensure all dependencies are properly checked

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- 🌐 [better-lyrics.boidu.dev](https://better-lyrics.boidu.dev) - Main project website
- 📚 [Whisper Documentation](https://github.com/openai/whisper)
- 📺 [yt-dlp Documentation](https://github.com/yt-dlp/yt-dlp)

---

<div align="center">
  <p>Made with ❤️ by the <a href="https://better-lyrics.boidu.dev" style="color: #FF0000; font-weight: bold;">better-lyrics</a> team</p>
  <p>⭐ Star us on GitHub if you find this useful!</p>
</div>