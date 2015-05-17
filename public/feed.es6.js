export class RssFeed {
  constructor() {
    this.items = [{
      title: "Some interesting item",
      GUID: "key",
      Channel: []
    }];
  }
  get() {
    return this.items;
  }
}
