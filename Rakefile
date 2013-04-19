require 'rake/clean'
require 'rbconfig'

# --- BUILD CONFIGURATION ---
UPX_EXE = 'C:/Apps/upx/bin/upx.exe'
S7ZIP_EXE = 'C:/tools/7za.exe'
# ---------------------------

task :default => :all

BUILD = 'build'
PKG = File.expand_path('pkg')

CLEAN.include(BUILD)
CLOBBER.include(PKG)


def dev_null
  if RbConfig::CONFIG['host_os'] =~ /mingw|mswin/
    'NUL'
  else
    '/dev/null'
  end
end

desc 'build all OS/arch flavors'
task :all => ['build:windows_386','build:linux_386','build:darwin_386']

namespace :all do
  desc 'build and shrink all exes'
  task :shrink => [:all] do
    Dir.chdir BUILD do
      Dir.glob('*').each do |d|
        Dir.chdir d do
          Dir.glob('uru*').each do |f|
            puts "---> upx shrinking #{d} #{f}"
            system "#{UPX_EXE} -9 #{f} > #{dev_null} 2>&1"
          end
        end
      end
    end
  end
end

namespace :build do
  %w[windows:386:0 linux:386:0 darwin:386:0].each do |tgt|
    os, arch, cgo = tgt.split(':')
    ext = (os == 'windows' ? '.exe' : '')

    desc "build #{os}/#{arch}"
    task :"#{os}_#{arch}" do |t|
      puts "---> building uru #{os}_#{arch} flavor"
      ENV['GOARCH'] = arch
      ENV['GOOS'] = os
      ENV['CGO_ENABLED'] = cgo
      system "go build -o #{BUILD}/#{t.name.split(':')[-1]}/uru_rt#{ext}"
    end
  end
end

desc 'archive all built exes'
task :package => 'package:all'

directory PKG
namespace :package do
  task :all => ['all:shrink',PKG] do
    ts = Time.now.strftime('%Y%m%dT%H%M')
    Dir.chdir BUILD do
      Dir.glob('*').each do |d|
        case d
        when /\A(darwin|linux)/
          puts "---> packaging #{d}"
          system "#{S7ZIP_EXE} a -tgzip -mx9 uru-#{$1}-#{ts}-bin-x86.gz ./#{d}/*  > #{dev_null} 2>&1"
          mv "uru-#{$1}-#{ts}-bin-x86.gz", PKG, :verbose => false
        when /\Awindows/
          puts "---> packaging #{d}"
          system "#{S7ZIP_EXE} a -t7z -mx9 uru-windows-#{ts}-bin-x86.7z ./#{d}/* > #{dev_null} 2>&1"
          mv "uru-windows-#{ts}-bin-x86.7z", PKG, :verbose => false
        end
      end
    end
  end
end
