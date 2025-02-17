#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "cJSON.h"
#include "cJSON.c"

#define MAX_LINE_LENGTH 1024
#define MAX_TOKEN_LENGTH 256

typedef struct
{
    char keyword[MAX_TOKEN_LENGTH];
    int count;
} KeywordCount;

typedef struct
{
    KeywordCount reserved_words[100];
    int reserved_count;
    KeywordCount functions[100];
    int function_count;
    KeywordCount variables[100];
    int variable_count;
    KeywordCount includes[100];
    int include_count;
} SuggestedConstraints;

typedef struct {
    char *name;
    KeywordCount *items;
    int count;
} ConstraintField;

const char *reserved_words[] = {
    "int", "float", "double", "char", "void", "if", "else", "while", "for", "return",
    "switch", "case", "break", "continue", "struct", "typedef", "const", "static", "unsigned", "signed"};

const int reserved_word_count = sizeof(reserved_words) / sizeof(reserved_words[0]);

int is_reserved_word(const char *word)
{
    for (int i = 0; i < reserved_word_count; i++)
    {
        if (strcmp(reserved_words[i], word) == 0)
        {
            return 1;
        }
    }
    return 0;
}

void add_to_constraints(KeywordCount *category, int *count, const char *keyword)
{
    for (int i = 0; i < *count; i++)
    {
        if (strcmp(category[i].keyword, keyword) == 0)
        {
            category[i].count++;
            return;
        }
    }
    strcpy(category[*count].keyword, keyword);
    category[*count].count = 1;
    (*count)++;
}

char *skip_string_literal(char *str)
{
    str++;
    while (*str && *str != '"')
    {
        if (*str == '\\')
            str++;
        if (*str)
            str++;
    }
    return *str ? str + 1 : str;
}

void analyze_code(FILE *file, SuggestedConstraints *constraints)
{
    char line[MAX_LINE_LENGTH];
    char prev_token[MAX_TOKEN_LENGTH] = "";
    int in_variable_declaration = 0;

    while (fgets(line, sizeof(line), file))
    {
        char *ptr = line;
        char token[MAX_TOKEN_LENGTH];

        while (*ptr)
        {
            while (*ptr && isspace(*ptr))
                ptr++;
            if (!*ptr)
                break;

            if (strncmp(ptr, "#include", 8) == 0)
            {
                ptr += 8;
                while (*ptr && isspace(*ptr))
                    ptr++;
                if (*ptr == '<' || *ptr == '"')
                {
                    char *start = ++ptr;
                    while (*ptr && *ptr != '>' && *ptr != '"')
                        ptr++;
                    if (*ptr)
                    {
                        *ptr = '\0';
                        add_to_constraints(constraints->includes, &constraints->include_count, start);
                    }
                }
                break;
            }

            if (*ptr == '"')
            {
                ptr = skip_string_literal(ptr);
                continue;
            }

            if (isalpha(*ptr) || *ptr == '_')
            {
                char *start = ptr;
                while (isalnum(*ptr) || *ptr == '_')
                    ptr++;
                int len = ptr - start;
                if (len < MAX_TOKEN_LENGTH)
                {
                    strncpy(token, start, len);
                    token[len] = '\0';

                    if (is_reserved_word(token))
                    {
                        add_to_constraints(constraints->reserved_words, &constraints->reserved_count, token);
                        if (strcmp(token, "int") == 0 || strcmp(token, "float") == 0 ||
                            strcmp(token, "char") == 0 || strcmp(token, "double") == 0)
                        {
                            in_variable_declaration = 1;
                        }
                    }
                    else
                    {
                        char *next = ptr;
                        while (*next && isspace(*next))
                            next++;
                        if (*next == '(')
                        {
                            add_to_constraints(constraints->functions, &constraints->function_count, token);
                            in_variable_declaration = 0;
                        }
                        else if (in_variable_declaration)
                        {
                            add_to_constraints(constraints->variables, &constraints->variable_count, token);
                        }
                    }
                }
            }
            else
            {
                if (*ptr == ';')
                {
                    in_variable_declaration = 0;
                }
                ptr++;
            }
        }
    }
}

cJSON* generate_constraints_json(SuggestedConstraints *constraints) { 
    cJSON *root = cJSON_CreateObject();
    cJSON_AddStringToObject(root, "status", "success");
    cJSON_AddStringToObject(root, "message", "Analysis completed successfully.");

    cJSON *data = cJSON_CreateObject();

    ConstraintField fields[] = {
        {"reserved_words", constraints->reserved_words, constraints->reserved_count},
        {"functions", constraints->functions, constraints->function_count},
        {"variables", constraints->variables, constraints->variable_count},
        {"includes", constraints->includes, constraints->include_count}
    };

    int num_fields = sizeof(fields) / sizeof(fields[0]);

    for (int i = 0; i < num_fields; i++) {
        cJSON *array = cJSON_AddArrayToObject(data, fields[i].name);
        for (int j = 0; j < fields[i].count; j++) {
            cJSON *item = cJSON_CreateObject();
            cJSON_AddStringToObject(item, "keyword", fields[i].items[j].keyword);
            cJSON_AddNumberToObject(item, "limit", fields[i].items[j].count);
            cJSON_AddItemToArray(array, item);
        }
    }

    cJSON_AddItemToObject(root, "data", data);
    
    return root;
}

// Find a keyword in the array and return the item or NULL if not found
cJSON* find_keyword_in_array(cJSON *array, const char *keyword) {
    if (!cJSON_IsArray(array)) return NULL;
    
    int size = cJSON_GetArraySize(array);
    for (int i = 0; i < size; i++) {
        cJSON *item = cJSON_GetArrayItem(array, i);
        cJSON *kw = cJSON_GetObjectItemCaseSensitive(item, "keyword");
        if (cJSON_IsString(kw) && strcmp(kw->valuestring, keyword) == 0) {
            return item;
        }
    }
    return NULL;
}

// Implementation of the check_constraints function similar to the Python version
cJSON* check_constraints(cJSON *user_kw_list, cJSON *kw_list) {
    cJSON *result = cJSON_CreateObject();
    int all_passed = 1;
    
    // For each category in kw_list
    cJSON *category, *user_category;
    const char *category_names[] = {"reserved_words", "functions", "variables", "includes"};
    int num_categories = sizeof(category_names) / sizeof(category_names[0]);
    
    for (int i = 0; i < num_categories; i++) {
        const char *category_name = category_names[i];
        cJSON *con_list = cJSON_GetObjectItemCaseSensitive(kw_list, category_name);
        
        // Skip if category doesn't exist
        if (!cJSON_IsArray(con_list)) continue;
        
        // Get corresponding user category
        cJSON *user_category = cJSON_GetObjectItemCaseSensitive(user_kw_list, category_name);
        if (!cJSON_IsArray(user_category)) {
            user_category = cJSON_CreateArray();
            cJSON_AddItemToObject(user_kw_list, category_name, user_category);
        }
        
        // Check each constraint in this category
        int array_size = cJSON_GetArraySize(con_list);
        for (int j = 0; j < array_size; j++) {
            cJSON *con = cJSON_GetArrayItem(con_list, j);
            
            // Add is_passed field to constraint
            cJSON_AddBoolToObject(con, "is_passed", 1);
            
            // Get constraint properties
            cJSON *con_keyword = cJSON_GetObjectItemCaseSensitive(con, "keyword");
            cJSON *con_type = cJSON_GetObjectItemCaseSensitive(con, "type");
            cJSON *con_limit = cJSON_GetObjectItemCaseSensitive(con, "limit");
            cJSON *con_active = cJSON_GetObjectItemCaseSensitive(con, "active");
            
            if (!cJSON_IsString(con_keyword) || !cJSON_IsString(con_type) || !cJSON_IsNumber(con_limit) || !cJSON_IsBool(con_active) ) {
                continue;
            }
            
            // Find corresponding user keyword
            cJSON *user_kw = find_keyword_in_array(user_category, con_keyword->valuestring);
            
            if (cJSON_IsTrue(con_active)){
                if (strcmp(con_type->valuestring, "eq") == 0) {
                    if (user_kw == NULL || cJSON_GetObjectItemCaseSensitive(user_kw, "limit")->valueint != con_limit->valueint) {
                        cJSON_ReplaceItemInObject(con, "is_passed", cJSON_CreateBool(0));
                        all_passed = 0;
                    }
                } else if (user_kw != NULL) {
                    int user_limit = cJSON_GetObjectItemCaseSensitive(user_kw, "limit")->valueint;
                    
                    if (strcmp(con_type->valuestring, "me") == 0 && user_limit < con_limit->valueint) {
                        cJSON_ReplaceItemInObject(con, "is_passed", cJSON_CreateBool(0));
                        all_passed = 0;
                    } else if (strcmp(con_type->valuestring, "le") == 0 && user_limit > con_limit->valueint) {
                        cJSON_ReplaceItemInObject(con, "is_passed", cJSON_CreateBool(0));
                        all_passed = 0;
                    } else if (strcmp(con_type->valuestring, "na") == 0 && user_limit > 0) {
                        cJSON_ReplaceItemInObject(con, "is_passed", cJSON_CreateBool(0));
                        all_passed = 0;
                    }
                }
            }
        }
    }
    
    cJSON_AddStringToObject(result, "status", all_passed ? "passed" : "failed");
    cJSON_AddItemToObject(result, "keyword_constraint", cJSON_Duplicate(kw_list, 1));
    
    return result;
}

int main(int argc, char *argv[]) {
    if (argc < 3) {
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Usage: program <source_file> <constraints_file>");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }

    // Open source file for analysis
    FILE *file = fopen(argv[1], "r");
    if (!file) {
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Error opening source file");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }

    // Analyze code and get user keyword list
    SuggestedConstraints constraints = {0};
    analyze_code(file, &constraints);
    fclose(file);
    
    // Fix: Properly get the data object from the generated constraints
    cJSON *generated = generate_constraints_json(&constraints);
    cJSON *user_kw_list = cJSON_GetObjectItemCaseSensitive(generated, "data");
    if (!user_kw_list) {
        cJSON_Delete(generated);
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Error generating user keyword list");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }
    
    // Read constraints file
    FILE *constraints_file = fopen(argv[2], "r");
    if (!constraints_file) {
        cJSON_Delete(generated);
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Error opening constraints file");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }
    
    // Read the constraints file content
    fseek(constraints_file, 0, SEEK_END);
    long file_size = ftell(constraints_file);
    fseek(constraints_file, 0, SEEK_SET);
    
    char *file_content = (char *)malloc(file_size + 1);
    if (!file_content) {
        fclose(constraints_file);
        cJSON_Delete(generated);
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Memory allocation failed");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }
    
    fread(file_content, 1, file_size, constraints_file);
    file_content[file_size] = '\0';
    fclose(constraints_file);
    
    // Parse constraints JSON
    cJSON *constraints_json = cJSON_Parse(file_content);
    free(file_content);
    
    if (!constraints_json) {
        cJSON_Delete(generated);
        cJSON *error = cJSON_CreateObject();
        cJSON_AddStringToObject(error, "status", "error");
        cJSON_AddStringToObject(error, "message", "Error parsing constraints JSON");
        char *json_string = cJSON_Print(error);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(error);
        return 1;
    }
    
    // Check constraints
    cJSON *result = check_constraints(user_kw_list, constraints_json);
    
    // Print result
    char *json_string = cJSON_Print(result);
    if (json_string) {
        printf("%s\n", json_string);
        free(json_string);
    }
    
    // Cleanup
    cJSON_Delete(result);
    cJSON_Delete(constraints_json);
    cJSON_Delete(generated);  // Delete the entire generated object
    
    return 0;
}