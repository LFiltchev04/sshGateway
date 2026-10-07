package internalrouter


//for transferring data from the frontend to the backend
type BufferEntry struct {
	channel []byte
	req []byte
}

