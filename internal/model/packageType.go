package model

type PackageType string

const (
	WithoutPackage PackageType = "noPackage"
	PlasticBag     PackageType = "plasticBag"
	Box            PackageType = "box"
	Film           PackageType = "film"
)
