package model

type Package struct {
	Package          PackageType `db:"package"`
	PackageSurcharge int         `db:"surcharge"`
	PackageMaxWeight int         `db:"max_weight"`
}
