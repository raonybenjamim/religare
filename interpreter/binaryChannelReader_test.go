package interpreter

import (
	"lazarus/models"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

var mockChecksum = createMockChecksum(models.ChecksumBits)
var mockMessageSize = zfill(strconv.FormatInt(20, 2), models.MessageSizeBits)
var validHeadersForTextMessage = models.MessageType.Text + mockChecksum + mockMessageSize

func TestGetHeaders(t *testing.T) {
	communicationChannel := make(chan models.Binary, len(validHeadersForTextMessage))

	channelReader := ChannelReader{
		Channel: communicationChannel,
	}

	loadChannel(communicationChannel, validHeadersForTextMessage)

	headers, err := channelReader.readHeadersFromChannel()

	// Check function should return true
	if err != nil {
		t.Fatalf("Error while reading headers %v", err)
	}

	if headers.Checksum != mockChecksum {
		t.Errorf("Header Checksum was not equal. Expected: %v, Got: %v", mockChecksum, headers.Checksum)
	}

	if headers.MessageSizeBytes != 20 {
		t.Errorf("Header Message Size was not equal. Expected: %v, Got: %v", mockMessageSize, headers.MessageSizeBytes)
	}

	if headers.MessageType != models.MessageType.Text {
		t.Errorf("Header Message Type was not equal. Expected %v, Got: %v", models.MessageType.Text, headers.MessageType)
	}
}

func loadChannel(ch chan models.Binary, value string) {
	for _, bit := range value {
		switch bit {
		case '0':
			ch <- models.Zero
		case '1':
			ch <- models.One
		}
	}
}

func createMockChecksum(size int) string {
	checksum := make([]byte, size) // Create a byte slice of the given size

	// Generate random 0s and 1s
	for i := 0; i < size; i++ {
		// Generate a random number between 0 and 1
		bit := rand.Intn(2)

		// Convert the random number to a byte (0 or 1) and store it in the checksum
		checksum[i] = byte('0' + bit)
	}

	return string(checksum) // Convert byte slice to string and return
}

func zfill(str string, length int) string {
	if len(str) >= length {
		return str // If the string is already equal or longer than the desired length, return it as is
	}
	zerosNeeded := length - len(str)          // Calculate the number of zeros needed
	zeros := strings.Repeat("0", zerosNeeded) // Create a string of zeros
	return zeros + str                        // Concatenate the zeros with the original string
}
