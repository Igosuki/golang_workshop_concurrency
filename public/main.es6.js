import {
  bootstrap,
  ComponentAnnotation as Component,
  ViewAnnotation as View,
  For,
  EventEmitter,
  DirectiveAnnotation as Directive,
} from './angular2/angular2';

import {name} from './exports.es6';

import {RssFeed} from './feed.es6';

@Component({
  selector: 'auto-reloader'
})
class AutoReloader {
  constructor(@EventEmitter('reload') reload:Function) {
    setInterval(() => reload(Date.now()), 10000)
  }
}

@Component({
  selector: 'colorwheel'
})
@View({
  template: `
    <button (click)="nextColor()" [style.background-color]="currentColor">Click The Wheel</button>
  `
})
class ColorWheel {
  constructor() {
    this.colors = ["green", "red", "yellow", "blue"]
    this.colorIdx = -1;
    this.nextColor()
  }
  nextColor() {
    console.log("Next color");
    this.colorIdx += 1;
    if(this.colorIdx >= this.colors.length) {
      this.colorIdx = 0;
    }
    this.currentColor = this.colors[this.colorIdx]
  }
}

@Component({
  selector: 'rss-app',
  injectables: [RssFeed]
})
@View({
  templateUrl: 'templates/feed.html',
  directives: [For, AutoReloader, ColorWheel]
})
class RssApp {
  constructor(rssFeed:RssFeed){
    this.items = rssFeed.get();
  }
  read(index) {
    this.items.splice(index, 1);
  }
  reloadFeed(time) {
    console.log(time);
    this.items = rssFeed.get();
  }
}

export function run () {
  bootstrap(RssApp);
}
