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

int main(int argc, char *argv[]) {
    cJSON *root = cJSON_CreateObject();
    
    if (argc < 2) {
        cJSON_AddStringToObject(root, "status", "error");
        cJSON_AddStringToObject(root, "message", "No input file specified");
        char *json_string = cJSON_Print(root);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(root);
        return 1;
    }

    FILE *file = fopen(argv[1], "r");
    if (!file) {
        cJSON_AddStringToObject(root, "status", "error");
        cJSON_AddStringToObject(root, "message", "Error opening input file");
        char *json_string = cJSON_Print(root);
        if (json_string) {
            printf("%s\n", json_string);
            free(json_string);
        }
        cJSON_Delete(root);
        return 1;
    }

    SuggestedConstraints constraints = {0};
    analyze_code(file, &constraints);
    fclose(file);

    cJSON *json = generate_constraints_json(&constraints);
    char *json_string = cJSON_Print(json);

    if (json_string) {
        printf("%s\n", json_string);
        free(json_string);
    }

    cJSON_Delete(json);
    return 0;
}