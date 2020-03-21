package wordsmith

type ManualWord struct {
	SingularStr,
	LowercaseAbbreviationStr,
	AbbreviationStr,
	PluralStr,
	RouteNameStr,
	KebabNameStr,
	PluralRouteNameStr,
	UnexportedVarNameStr,
	PluralUnexportedVarNameStr,
	PackageNameStr,
	SingularPackageNameStr,
	SingularCommonNameStr,
	ProperSingularCommonNameWithPrefixStr,
	PluralCommonNameStr,
	SingularCommonNameWithPrefixStr,
	PluralCommonNameWithPrefixStr string
}

func (mw *ManualWord) Singular() string {
	return mw.SingularStr
}

func (mw *ManualWord) LowercaseAbbreviation() string {
	return mw.AbbreviationStr
}

func (mw *ManualWord) Abbreviation() string {
	return mw.AbbreviationStr
}

func (mw *ManualWord) Plural() string {
	return mw.PluralStr
}

func (mw *ManualWord) RouteName() string {
	return mw.RouteNameStr
}

func (mw *ManualWord) KebabName() string {
	return mw.KebabNameStr
}

func (mw *ManualWord) PluralRouteName() string {
	return mw.PluralRouteNameStr
}

func (mw *ManualWord) UnexportedVarName() string {
	return mw.UnexportedVarNameStr
}

func (mw *ManualWord) PluralUnexportedVarName() string {
	return mw.PluralUnexportedVarNameStr
}

func (mw *ManualWord) PackageName() string {
	return mw.PackageNameStr
}

func (mw *ManualWord) SingularPackageName() string {
	return mw.SingularPackageNameStr
}

func (mw *ManualWord) SingularCommonName() string {
	return mw.SingularCommonNameStr
}

func (mw *ManualWord) ProperSingularCommonNameWithPrefix() string {
	return mw.ProperSingularCommonNameWithPrefixStr
}

func (mw *ManualWord) PluralCommonName() string {
	return mw.PluralCommonNameStr
}

func (mw *ManualWord) SingularCommonNameWithPrefix() string {
	return mw.SingularCommonNameWithPrefixStr
}

func (mw *ManualWord) PluralCommonNameWithPrefix() string {
	return mw.PluralCommonNameWithPrefixStr
}
