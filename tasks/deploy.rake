require 'erb'

if DEPLOY_MODE
  namespace :deploy do
    desc 'deploy uru binaries to sourceforge.net'
    task :sf => ['package:all'] do
      # upload via sftp/psftp and set files as new default with Net::HTTP::Put.new
      #   http://sourceforge.net/p/forge/community-docs/Using%20the%20Release%20API/
      #   http://www.ruby-doc.org/stdlib-2.0.0/libdoc/net/http/rdoc/Net/HTTP.html
      #   http://www.rubyinside.com/nethttp-cheat-sheet-2940.html
      #   https://blogs.oracle.com/edwingo/entry/ruby_multipart_post_put_request
      SFDeployer.deploy_files
      SFDeployer.set_default_files
    end
  end


  # helpers
  module SFDeployer

    @windows_archive = File.expand_path("pkg/uru-#{VER}-windows-#{CPU}.7z")
    @linux_archive = File.expand_path("pkg/uru-#{VER}-linux-#{CPU}.tar.gz")
    @darwin_archive = File.expand_path("pkg/uru-#{VER}-darwin-#{CPU}.tar.gz")

    def self.sftp_batch_script
      data = {
        :version => VER,
        :windows_archive => @windows_archive,
        :linux_archive => @linux_archive,
        :darwin_archive => @darwin_archive
      }

      erb = ERB.new(UruDeployConfig.sftp_script_template)
      erb.result(binding)
    end

    def self.deploy_files
      batch_file = 'sftp_batch_file'
      File.open(batch_file, 'w') do |f|
        f.write(sftp_batch_script)
	  end

      begin
        command = "#{SFTP_EXE} -i #{UruDeployConfig.sf_private_key} #{UruDeployConfig.sf_user}@frs.sourceforge.net -b #{batch_file}"
        puts "---> deploying to urubinaries at sf.net"
        system "#{command} > #{dev_null} 2>&1"
      ensure
        File.unlink(batch_file)
      end
    end

    def self.set_default_files
      puts '---> setting default files'
    end

  end
end
