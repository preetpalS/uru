desc 'test all uru packages'
task :test => 'test:all'

namespace :test do
  task :all => ['test:env']

  task :env do
    puts "---> testing `env` package"
    system "go test #{GO_PKG_ROOT}/env"
  end
end
