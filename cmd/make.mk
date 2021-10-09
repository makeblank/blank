VPATH ?= $(BLANK_PATH)
target = $(firstword $(MAKECMDGOALS))
others = $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

export VPATH

ifdef target
paths  = $(subst :, ,$(VPATH))
incs   = $(paths:%=-I %)
files  = $(paths:%=%/$(target).mk)
found  = $(realpath $(files))
found := $(firstword $(found))

ifdef found
$(target):
	@$(MAKE) --no-silent --no-print-directory $(incs) -f $(found) $(others)
else
$(target):
	$(error No target found: $@)
endif
endif

.DEFAULT: ;
