import {
  bootstrap,
  ComponentAnnotation as Component,
  ViewAnnotation as View,
  For
} from './angular2/angular2';

import {name} from './exports.es6';

import {RssFeed} from './feed.es6'

@Component({
  selector: 'rss-app',
  injectables: [RssFeed]
})
@View({
  templateUrl: 'templates/rss-feed.html',
  directives: [For]
})
class RssApp {
  constructor(rssFeed:RssFeed){
    this.items = rssFeed.get();
  }
  read(index) {
    this.items.splice(index, 1);
  }
}

export function run () {
  bootstrap(RssApp);
}
