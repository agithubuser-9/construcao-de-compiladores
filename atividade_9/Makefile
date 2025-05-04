ASSEMBLY_DIR := assembly
BIN_DIR := bin
RUNTIME := assembly/runtime.s


SOURCES := $(shell find $(ASSEMBLY_DIR) -maxdepth 1 -name "*.s" ! -name "runtime.s")
TARGETS := $(foreach src,$(SOURCES),$(BIN_DIR)/$(notdir $(basename $(src))))

all: $(TARGETS)

# Regras por arquivo .s
define BUILD_RULE
$(BIN_DIR)/$(notdir $(basename $(1))): $(1) $(RUNTIME)
	@mkdir -p $(BIN_DIR)
	as -o tmp_$$(notdir $(basename $(1))).o $(1)
	as -o tmp_runtime.o $(RUNTIME)
	ld -o $$@ tmp_$$(notdir $(basename $(1))).o tmp_runtime.o
	rm -f tmp_$$(notdir $(basename $(1))).o tmp_runtime.o
	@echo "✅ Gerado: $$@"
endef

$(foreach src,$(SOURCES),$(eval $(call BUILD_RULE,$(src))))

clean:
	rm -rf $(BIN_DIR)/*

test:
	@for exe in $(wildcard $(BIN_DIR)/*); do \
		nome=$$(basename $$exe); \
		saida=$$(./$$exe); \
		esperado=$$(cat resultados/$$nome.txt 2>/dev/null); \
		if [ "$$saida" = "$$esperado" ]; then \
			echo "✅ $$nome passou! Resultado: $$saida"; \
		else \
			echo "❌ $$nome falhou! Esperado: $$esperado, Obtido: $$saida"; \
		fi \
	done

