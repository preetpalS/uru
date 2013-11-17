if DEPLOY_MODE

  namespace :deploy do
    desc 'deploy archives to urubinaries at sourceforge.net'
    task :sf do
      # upload via pscp and set files as new default with Net::HTTP::Put.new
      #   http://sourceforge.net/p/forge/community-docs/Using%20the%20Release%20API/
      #   http://www.ruby-doc.org/stdlib-2.0.0/libdoc/net/http/rdoc/Net/HTTP.html
      #   http://www.rubyinside.com/nethttp-cheat-sheet-2940.html
      #   https://blogs.oracle.com/edwingo/entry/ruby_multipart_post_put_request
      puts '---> deploying to sf.net using:'
      system "#{SCP_EXE} --version"
    end
  end
end
