'use strict';

var gulp = require('gulp');
var del = require('del');
var Builder = require('systemjs-builder');
var traceur = require('gulp-traceur');

var http = require('http');
var connect = require('connect');
var serveStatic = require('serve-static');
var open = require('open');

gulp.task('clean', function (done) {
  del(['dev'], done);
});

gulp.task('build:angular2', function () {
  var builder = new Builder({
    paths: {
      'angular2/*': 'node_modules/angular2/es6/dev/*.es6',
      rx: 'node_modules/angular2/node_modules/rx/dist/rx.js'
    },
    meta: {
      rx: {
        format: 'cjs'
      }
    }
  });
  return builder.build('angular2/angular2', './public/lib/angular2.js', {});
});

gulp.task('build:lib', ['build:angular2'], function () {
  gulp.src([
    './node_modules/gulp-traceur/node_modules/traceur/bin/traceur-runtime.js',
    './node_modules/gulp-traceur/node_modules/traceur/bin/traceur.js',
    './node_modules/angular2/node_modules/zone.js/dist/zone.js',
    './node_modules/es6-module-loader/dist/es6-module-loader-sans-promises.js',
    './node_modules/es6-module-loader/dist/es6-module-loader-sans-promises.js.map',
    './node_modules/reflect-metadata/Reflect.js',
    './node_modules/reflect-metadata/Reflect.js.map',
    './node_modules/systemjs/dist/system.js',
    './node_modules/systemjs/dist/system.js.map'
  ])
  .pipe(gulp.dest('./public/lib'));
});

gulp.task('build', function () {
  return gulp.src('./js/**/*.js')
    .pipe(traceur({
      sourceMaps: 'inline',
      modules: 'instantiate',
      annotations: true,
      memberVariables: true,
      types: true
    }))
    .pipe(gulp.dest('./dev'));
});

gulp.task('serve', ['build:lib', 'build'], function () {
  var port = 5555
  var app;

  gulp.watch('./js/**', ['build']);

  app = connect().use(serveStatic(__dirname));
  http.createServer(app).listen(port, function () {
    open('http://localhost:' + port);
  });
});

