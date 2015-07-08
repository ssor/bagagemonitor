package controllers

import ()

type ImageInfo struct {
	Name    string
	Deleted bool
}

func (this *ImageInfo) SetDeleted() {
	this.Deleted = true
}

type ImageInfoList []*ImageInfo

func (this ImageInfoList) Contains(name string) bool {
	return this.Find(name) != nil
}
func (this ImageInfoList) Find(name string) *ImageInfo {
	for _, ii := range this {
		if ii.Name == name {
			return ii
		}
	}
	return nil
}
func (this ImageInfoList) RegisterImage(name string) ImageInfoList {
	if this.Contains(name) {
		return this
	}
	return append(this, &ImageInfo{name, false})
}
func (this ImageInfoList) Clear() (list ImageInfoList) {
	for _, ii := range this {
		if ii.Deleted == false {
			list = append(list, ii)
		}
	}
	return
}
