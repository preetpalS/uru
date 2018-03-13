# Author: Jon Maken, All Rights Reserved
# Licence: 3-clause BSD

require 'digest/sha2'
require 'erb'
require 'fileutils'

namespace :scoop do
  scoop_root = File.join(BUILD, 'scoop')
  archive_path = "#{PKG}/uru-#{VER}-windows-x86.7z"

  task :prep do
    unless File.exist?(archive_path)
      abort "---> FAILED to find `pkg/uru-#{VER}-windows-x86.7z` needed to build scoop manifest"
    end

    FileUtils.mkdir_p(scoop_root) unless Dir.exist?(scoop_root)
  end

  task :templates do
    template_root = File.join(File.dirname(__FILE__), 'templates')
    archive_sha256 = Digest::SHA256.file(archive_path).hexdigest

    data = {
      url: "https://bitbucket.org/jonforums/uru/downloads/uru-#{VER}-windows-x86.7z",
      ver: "#{VER}",
      chksum: archive_sha256,
    }

    %w[uru.json.erb].each do |t|
      File.open(File.join(scoop_root, t.gsub('.erb','')) ,'w+') do |f|
        f.write(ERB.new(File.read(File.join(template_root, t)), nil, '<>').result(binding))
      end
    end
  end

  desc 'build uru.json scoop app manifest'
  task :package => [:prep, :templates] do
    Dir.chdir(scoop_root) do
      FileUtils.cp("uru.json", PKG, :preserve => false)
    end
  end

  task :deploy do
    abort "---> TODO implement `scoop:deploy` task"
  end
end
