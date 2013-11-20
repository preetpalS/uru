# The deployment task requires two non-revision controlled files in the root
# of the project directory: deploy_config.rb and cacert.pem.
#
# a cacert.pem can be obtained from http://curl.haxx.se/docs/caextract.html
#
# deploy_config.rb looks similar to:
#
#  UruDeployConfig = OpenStruct.new(
#    :sf_user => 'YOUR SF.NET USERNAME',
#    :sf_private_key => '/path/to/your/private/key/registered/at/sf.net',
#    :sf_api_key => 'YOUR SF.NET API KEY',
#  )
#
#  UruDeployConfig.sftp_script_template =<<-EOF
#  cd /home/frs/project/urubinaries/uru
#  mkdir <%= data[:version] %>
#  cd <%= data[:version] %>
#  put <%= data[:windows_archive] %>
#  put <%= data[:linux_archive] %>
#  put <%= data[:darwin_archive] %>
#  EOF

if DEPLOY_MODE

  require 'erb'
  require 'net/https'

  namespace :deploy do
    desc 'deploy uru binaries to sourceforge.net'
    task :sf => ['package:all'] do
      # TODO likely these two are too quick for sf.net's infrastructure; split
      SFDeployer.deploy_files
      SFDeployer.set_default_files
    end
  end


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

    # TODO implement error handling
    def self.set_default_files
      # http://sourceforge.net/p/forge/community-docs/Using%20the%20Release%20API/
      http = Net::HTTP.new('sourceforge.net', 443)
      http.use_ssl = true
      http.verify_mode = ::OpenSSL::SSL::VERIFY_PEER

      store = ::OpenSSL::X509::Store.new
      store.add_file('cacert.pem')
      http.cert_store = store

      {
        File.basename(@windows_archive) => %w[windows],
        File.basename(@linux_archive) => %w[linux bsd solaris others],
        File.basename(@darwin_archive) => %w[mac]
      }.each do |f, d|
        puts "---> setting #{f} as default for #{d.join(', ')}"
        req = Net::HTTP::Put.new("/projects/urubinaries/files/uru/#{VER}/#{f}")
        req['Accept'] = 'application/json'
        req.set_form_data 'api_key' => UruDeployConfig.sf_api_key,
                          'default' => (d.size == 1 ? d[0] : d)
        # TODO parse JSON response for error that looks similar to
        #        { "error": "blah blah blah" }
        res = http.request(req)
      end
    end

  end
end
