package frontendsrc

const (
	queryFilterDotTSFile = `const pageKey = "page";
const limitKey = "limit";
const createdBeforeKey = "createdBefore";
const createdAfterKey = "createdAfter";
const updatedBeforeKey = "updatedBefore";
const updatedAfterKey = "updatedAfter";
const sortByKey = "sortBy";

const enum sortBys {
  ASCENDING = "asc",
  DESCENDING = "desc",
}

const defaultPage = 1;
const defaultLimit = 20;
const defaultCreatedBefore = 0;
const defaultCreatedAfter = 0;
const defaultUpdatedBefore = 0;
const defaultUpdatedAfter = 0;
const defaultSortBy = sortBys.ASCENDING;

export class QueryFilter {
  public page: number;
  public limit: number;
  public createdBefore: number;
  public createdAfter: number;
  public updatedBefore: number;
  public updatedAfter: number;
  public sortBy: string;

  modifyURL(u: URL): void {
    if (this.page !== 0) {
      u.searchParams.append(pageKey, this.page.toString());
    }
    if (this.limit !== 0) {
      u.searchParams.append(limitKey, this.limit.toString());
    }
    if (this.createdBefore !== 0) {
      u.searchParams.append(createdBeforeKey, this.createdBefore.toString());
    }
    if (this.createdAfter !== 0) {
      u.searchParams.append(createdAfterKey, this.createdAfter.toString());
    }
    if (this.updatedBefore !== 0) {
      u.searchParams.append(updatedBeforeKey, this.updatedBefore.toString());
    }
    if (this.updatedAfter !== 0) {
      u.searchParams.append(updatedAfterKey, this.updatedAfter.toString());
    }
    if (this.sortBy !== "") {
      u.searchParams.append(sortByKey, this.sortBy.toString());
    }
  }

  constructor(params: URLSearchParams) {
    this.page = params.get(pageKey)
      ? parseInt(params.get(pageKey) || defaultPage.toString())
      : defaultPage;
    this.limit = params.get(limitKey)
      ? parseInt(params.get(limitKey) || defaultLimit.toString())
      : defaultLimit;
    this.createdBefore = params.get(createdBeforeKey)
      ? parseInt(params.get(createdBeforeKey) || defaultCreatedBefore.toString())
      : defaultCreatedBefore;
    this.createdAfter = params.get(createdAfterKey)
      ? parseInt(params.get(createdAfterKey) || defaultCreatedAfter.toString())
      : defaultCreatedAfter;
    this.updatedBefore = params.get(updatedBeforeKey)
      ? parseInt(params.get(updatedBeforeKey) || defaultUpdatedBefore.toString())
      : defaultUpdatedBefore;
    this.updatedAfter = params.get(updatedAfterKey)
      ? parseInt(params.get(updatedAfterKey) || defaultUpdatedAfter.toString())
      : defaultUpdatedAfter;

    const sb = params.get(sortByKey);
    if (
      sb &&
      sb.toLowerCase() !== sortBys.ASCENDING &&
      sb.toLowerCase() !== sortBys.DESCENDING
    ) {
      this.sortBy = defaultSortBy;
    } else {
      this.sortBy = params.get(sortByKey) || defaultSortBy;
    }
  }
}
`
)

func queryFilterDotTS() string {
	return queryFilterDotTSFile
}
