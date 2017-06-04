package rec 

 /*
  #include <stdio.h>
  #include <unistd.h>
  #include <termios.h>
  char getch(){
      char ch = 0;
      struct termios old = {0};
      fflush(stdout);
      if( tcgetattr(0, &old) < 0 ) perror("tcsetattr()");
      old.c_lflag &= ~ICANON;
      old.c_lflag &= ~ECHO;
      old.c_cc[VMIN] = 1;
      old.c_cc[VTIME] = 0;
      if( tcsetattr(0, TCSANOW, &old) < 0 ) perror("tcsetattr ICANON");
      if( read(0, &ch,1) < 0 ) perror("read()");
      old.c_lflag |= ICANON;
      old.c_lflag |= ECHO;
      if(tcsetattr(0, TCSADRAIN, &old) < 0) perror("tcsetattr ~ICANON");
      return ch;
  }
 */
import "C"

import (
        "fmt"
        "github.com/gordonklaus/portaudio"
        wave "github.com/zenwerk/go-wave"
        "math/rand"
        "os"
)

func errCheck(err error) {
        if err != nil {
                panic(err)
        }
}

func randStringBytesRmndr(n int) string {
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    return string(b)
}

func Record(verbose bool) (audioFileName string, result bool) {
        audioFileName = randStringBytesRmndr(10) + ".wav"
        result = false

        if verbose {
                fmt.Println("Recording. Press ESC to quit.")
        }

        waveFile, err := os.Create(audioFileName)
        errCheck(err)

        portaudio.Initialize()

        inputChannels := 1
        outputChannels := 0
        sampleRate := 44100
        framesPerBuffer := make([]byte, 64)
        stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(framesPerBuffer), framesPerBuffer)
        errCheck(err)

        param := wave.WriterParam{
                Out:           waveFile,
                Channel:       inputChannels,
                SampleRate:    sampleRate,
                BitsPerSample: 8,
        }

        waveWriter, err := wave.NewWriter(param)
        errCheck(err)

        go func() {
                key := C.getch()
                fmt.Println()
                if key == 27 {
                        waveWriter.Close()
                        stream.Close()
                        portaudio.Terminate()
                        result = true

                        return
                }
        }()

        errCheck(stream.Start())
        for {
                errCheck(stream.Read())
                _, err := waveWriter.Write([]byte(framesPerBuffer)) // WriteSample16 for 16 bits
                errCheck(err)
        }
        errCheck(stream.Stop())

        return
}
