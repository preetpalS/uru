# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

desc 'test uru source files'
task :test => 'test:all'

namespace :test do
  task :all => ['build:prep','test:main','test:env','test:exec','test:command']

  task :main do
    puts "\n---> testing `main` package"
    system "go test #{GO_PKG_ROOT}"
  end

  task :env do
    puts "\n---> testing `env` package"
    system "go test #{GO_PKG_ROOT}/env"
  end

  task :exec do
    puts "\n---> testing `exec` package"
    system "go test #{GO_PKG_ROOT}/exec"
  end

  task :command do
    puts "\n---> testing `command` package"
    system "go test #{GO_PKG_ROOT}/command"
  end
end
