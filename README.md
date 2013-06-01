# Unleash Ruby

Uru is a lightweight, multi-platform command line tool that helps you use the
multiple rubies (currently MRI and JRuby) on your 32/64-bit Windows, Linux, or
OS X systems.

While there are a number of fantastic ruby environment managers such as [RVM][1],
[rbenv][2], [pik][3], and [chruby][4], none are truely multi-platform, and all
provide different levels of awesomeness. In contrast, [uru][5] is a micro-kernel.
It provides a core set of minimal complexity, cross-platform ruby use enhancers.

For the most part, Ruby is a multi-platform joy to use. Wouldn't it be great if
your ruby environment manager spoke multi-platform?

# Easy to Install

The quickest path to uru zen is to [download][download] the binary archive for
your platform type, extract its contents to a directory already on `PATH`, and
perform one of the following installs. To build and install uru from Go source,
or get more in-depth details, review the [installation and usage][usage] info.

## Windows systems

~~~ console
:: assuming C:\tools is on PATH and uru_rt.exe was extracted to C:\tools
C:\tools>uru_rt admin install

:: register a pre-existing "system" ruby already placed on PATH as part
:: of cmd.exe initialization
C:\tools>uru_rt admin add system
~~~

## Linux and OS X systems

~~~ console
# assuming ~/bin is on PATH and uru_rt was extracted to ~/bin
$ cd ~/bin && chmod +x uru_rt

# append to ~/.profile on Ubuntu, or to ~/.zshrc on Zsh
$ echo 'eval "$(uru_rt admin install)"' >> ~/.bash_profile

# register a pre-existing "system" ruby already placed on PATH from bash/Zsh
# startup configuration files
$ uru_rt admin add system

# restart shell
$ exec $SHELL -l
~~~

# Easy to Use

While examples of uru's core commands are [available here][examples], once your
installed rubies have been registered with uru, usage is similar to:

~~~ bash
$ uru ls
        173: jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Serv...
 =>  system: ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

$ uru 173
---> Now using jruby 1.7.3

$ uru ls
 =>     173: jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Serv...
     system: ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

$ jruby --version
jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]

$ uru sys
---> Now using ruby 2.1.0dev

$ uru ls
        173: jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Serv...
 =>  system: ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

$ ruby --version
ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

$ uru gem li rake
jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]

rake (10.0.3)

ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

rake (10.0.4, 0.9.6)

$ uru ruby hello.rb
jruby 1.7.3 (1.9.3p385) 2013-02-21 dac429b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]

hello there

ruby 2.1.0dev (2013-05-06 trunk 40593) [i686-linux]

hello there
~~~

[download]: https://bitbucket.org/jonforums/uru/wiki/Downloads
[usage]: https://bitbucket.org/jonforums/uru/wiki/Usage
[examples]: https://bitbucket.org/jonforums/uru/wiki/Examples

[1]: https://rvm.io/
[2]: https://github.com/sstephenson/rbenv
[3]: https://github.com/vertiginous/pik
[4]: https://github.com/postmodern/chruby
[5]: https://bitbucket.org/jonforums/uru
