package lib

/*
Formats Object
*/
type Formats []string

/*
Add Method
*/
func (fs *Formats) Add(f string) {
	*fs = append(*fs, f)
}

/*
Copy Method
*/
func (fs *Formats) Copy(fs2 Formats) {
	for _, f := range fs2 {
		fs.Add(f)
	}
}

/*
Size Method
*/
func (fs *Formats) Size() int {
	return len(*fs)
}

/*
FormatMixin Object
*/
type FormatMixin struct {
	formats *Formats
}

func (fsm *FormatMixin) initFormats() {
	fsm.formats = &Formats{}
}

/*
AddFormat Method
*/
func (fsm *FormatMixin) AddFormat(fs ...string) *FormatMixin {
	fsm.formats.Copy(fs)
	return fsm
}

/*
Formats Method
*/
func (fsm *FormatMixin) Formats() *Formats {
	return fsm.formats
}

/*
ClearFormats Method
*/
func (fsm *FormatMixin) ClearFormats() *FormatMixin {
	fsm.formats = &Formats{}
	return fsm
}

/*
CopyFormats Method
*/
func (fsm *FormatMixin) CopyFormats(fsm2 FormatMixin) {
	fsm.formats.Copy(*fsm2.Formats())
}
