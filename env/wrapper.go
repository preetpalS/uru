// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// +build !windows

package env

var BashWrapper = `uru()
{
  export URU_INVOKER='bash'

  # uru_rt must already be on PATH
  uru_rt "$@"

  if [[ -d "$URU_HOME" ]]; then
    if [[ -f "$URU_HOME/uru_lackee" ]]; then
      . "$URU_HOME/uru_lackee"
      export PATH="${NEW_PATH}"

	  if [[ "$UNSET_GEM_HOME" == "yes" ]]; then
	    unset GEM_HOME UNSET_GEM_HOME
	  else
        export GEM_HOME="${NEW_GEM_HOME}"
	  fi
    fi
  else
    if [[ -f "$HOME/.uru/uru_lackee" ]]; then
      . "$HOME/.uru/uru_lackee"
      export PATH="${NEW_PATH}"

	  if [[ "$UNSET_GEM_HOME" == "yes" ]]; then
	    unset GEM_HOME UNSET_GEM_HOME
	  else
        export GEM_HOME="${NEW_GEM_HOME}"
      fi
    fi
  fi
}
`
