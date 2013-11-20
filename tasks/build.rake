# Author: Jon Maken, All Rights Reserved
# License: 3-clause BSD

desc 'build all OS/arch flavors'
task :all => BUILDS

namespace :build do
  puts "\n  *** DEVELOPMENT build mode ***\n\n" if URU_OPTS[:devbuild]

  %W[windows:#{ARCH}:0 linux:#{ARCH}:0 darwin:#{ARCH}:0].each do |tgt|
    os, arch, cgo = tgt.split(':')
    ext = (os == 'windows' ? '.exe' : '')

    desc "build #{os}/#{arch}"
    task :"#{os}_#{arch}" do |t|
      puts "---> building uru #{os}_#{arch} flavor"
      ENV['GOARCH'] = arch
      ENV['GOOS'] = os
      ENV['CGO_ENABLED'] = cgo
      system %Q{go build -ldflags "-s" -o #{BUILD}/#{t.name.split(':')[-1]}/uru_rt#{ext}}
    end
  end
end
