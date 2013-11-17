# --- CUSTOMIZE BUILD CONFIGURATION ---
GO_PKG_ROOT = 'bitbucket.org/jonforums/uru'
S7ZIP_EXE = 'C:/tools/7za.exe'
SCP_EXE = 'C:/tools/pscp.exe'
# -------------------------------------

# load project archive deployment configuration file if present
begin
  require File.expand_path('deploy_config')
  DEPLOY_MODE = Module.constants.include?(:UruDeployConfig)
rescue LoadError
end

# load modularized rake tasks
Dir['tasks/*.rake'].sort.each { |f| load f }
