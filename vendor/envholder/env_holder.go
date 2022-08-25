package envholder

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// EnvHolder holds all variables from .env file into memory
type EnvHolder struct {
	vars []EnvVar
}

// EnvVar is a key => value pair representing one variable from the .env file
type EnvVar struct {
	name string
	val string
}

// Returns the value for a given variable loaded from .env file by name
func (x EnvHolder) GetVar(var_name string) (string){
	for _, envVariable := range(x.vars) {
		if envVariable.name == var_name {
			return envVariable.val
		}
	}
	log.Fatalf("Variable %s not found in .env file.", var_name)
	return ""
}

// loadEnv reads from the .env file and acts as a factory method for the EnvHolder struct
func LoadEnv() (EnvHolder){
	file, ferr := os.Open(".env")
	if ferr != nil {
		log.Fatalf("Could not find .env file, details: %s", ferr.Error())
	}

	scanner := bufio.NewScanner(file)
	counter := 0
	var_holder := EnvHolder{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		split_line := strings.Split(line, "=")
		if len(split_line) != 2 {
			log.Printf("Invalid format for line %d of file: \"%s\"", counter, line)
			continue
		}
		loaded_variable := EnvVar{
			name: split_line[0],
			val: split_line[1],
		}
		var_holder.vars = append(var_holder.vars, loaded_variable)
		counter++
	}
	return var_holder
}