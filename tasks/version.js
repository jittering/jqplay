module.exports = function (grunt) {
  grunt.loadNpmTasks('grunt-text-replace');
  grunt.config('replace', {
    'index.tmpl': {
      src: ['assets/index.tmpl'],
      dest: 'public/root/index.tmpl',
      replacements: [{
        from: '#{ TIMESTAMP }',
        to: '<%= grunt.option("ts") %>'
      }]
    },
    'worker': {
      src: ['public/bower_components/ace-builds/src-min-noconflict/worker-xquery.js'],
      dest: 'public/js/worker-xquery.js',
      replacements: [],
    }
  })
}
