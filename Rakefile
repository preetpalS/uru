require 'rake/clean'
require 'rbconfig'

# --- BUILD CONFIGURATION ---
GO_PKG_ROOT = 'bitbucket.org/jonforums/uru'
UPX_EXE = 'C:/Apps/upx/bin/upx.exe'
S7ZIP_EXE = 'C:/tools/7za.exe'
# ---------------------------

task :default => :all

args = ARGV.dup
opts = {}
opts[:shrink] = args.delete('--shrink')
opts[:devbuild] = args.delete('--dev-build')

VER = /AppVersion\s*=\s*\`(\d{1,2}\.\d{1,2}\.\d{1,2})(\.\w+)?/.match(File.read('env/ui.go')) do |m|
  if m[2] != nil then m[1] + m[2] else m[1] end
end || 'NA'

ARCH = ENV['GOARCH'] || '386'
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

builds = %W[build:windows_#{ARCH} build:linux_#{ARCH} build:darwin_#{ARCH}]
desc 'build all OS/arch flavors'
task :all => builds

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


desc 'archive all built exes'
task :package => 'package:all'

directory PKG
pkg_prereqs = [PKG]
if opts[:shrink]
  pkg_prereqs.unshift('all:shrink')
else
  pkg_prereqs = builds + pkg_prereqs
end

namespace :package do
  task :all => pkg_prereqs do
    ts = `git rev-list --abbrev-commit -1 HEAD`.chomp
    cpu = case ARCH
          when 'amd64'
            'x64'
          when '386'
            'x86'
          else
            'NA'
          end
    Dir.chdir BUILD do
      Dir.glob('*').each do |d|
        case d
        when /\A(darwin|linux)/
          puts "---> packaging #{d}"
          tar = "uru-#{VER}-#{$1}.tar"
          archive = if opts[:devbuild]
                      "uru-#{VER}-#{ts}-#{$1}-#{cpu}.tar.gz"
                    else
                      "uru-#{VER}-#{$1}-#{cpu}.tar.gz"
                    end

          system "#{S7ZIP_EXE} a -ttar #{tar} ./#{d}/* > #{dev_null} 2>&1"
          system "#{S7ZIP_EXE} a -tgzip -mx9 #{archive} #{tar} > #{dev_null} 2>&1"
          mv archive, PKG, :verbose => false
          rm tar, :verbose => false
        when /\Awindows/
          puts "---> packaging #{d}"
          archive = if opts[:devbuild]
                      "uru-#{VER}-#{ts}-windows-#{cpu}.7z"
                    else
                      "uru-#{VER}-windows-#{cpu}.7z"
                    end

          system "#{S7ZIP_EXE} a -t7z -mx9 #{archive} ./#{d}/* > #{dev_null} 2>&1"
          mv archive, PKG, :verbose => false
        end
      end
    end
  end
end


desc 'test all uru packages'
task :test => 'test:all'

namespace :test do
  task :all => ['test:env']

  task :env do
    puts "---> testing `env` package"
    system "go test #{GO_PKG_ROOT}/env"
  end
end
