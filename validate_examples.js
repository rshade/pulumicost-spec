#!/usr/bin/env node

const Ajv = require("ajv").default;
const addFormats = require("ajv-formats").default;
const fs = require("fs");
const path = require("path");

// Create ajv instance
const ajv = new Ajv({
  strict: false,
  allErrors: true,
  loadSchema: false,
  addUsedSchema: false,
});
addFormats(ajv);

// Load schema
const schemaPath = "schemas/pricing_spec.schema.json";
const schema = JSON.parse(fs.readFileSync(schemaPath, "utf8"));

// Remove problematic $schema reference for validation
const cleanSchema = { ...schema };
delete cleanSchema.$schema;

// Compile schema
let validate;
try {
  validate = ajv.compile(cleanSchema);
  console.log("✅ Schema compilation successful");
} catch (err) {
  console.error("❌ Schema compilation error:", err.message);
  process.exit(1);
}

// Validate examples
let allValid = true;

// Validate pricing spec examples
const examplesDir = "examples/specs";
if (fs.existsSync(examplesDir)) {
  const jsonFiles = fs
    .readdirSync(examplesDir)
    .filter((f) => f.endsWith(".json"));

  for (const file of jsonFiles) {
    const filePath = path.join(examplesDir, file);
    try {
      const data = JSON.parse(fs.readFileSync(filePath, "utf8"));
      const valid = validate(data);

      if (valid) {
        console.log(`✅ ${file} is valid`);
      } else {
        console.error(`❌ ${file} is invalid:`);
        console.error(ajv.errorsText(validate.errors));
        allValid = false;
      }
    } catch (err) {
      console.error(`❌ Error processing ${file}:`, err.message);
      allValid = false;
    }
  }
}

// Validate recommendations examples (JSON syntax only)
const recommendationsDir = "examples/recommendations";
if (fs.existsSync(recommendationsDir)) {
  const jsonFiles = fs
    .readdirSync(recommendationsDir)
    .filter((f) => f.endsWith(".json"));

  for (const file of jsonFiles) {
    const filePath = path.join(recommendationsDir, file);
    try {
      JSON.parse(fs.readFileSync(filePath, "utf8"));
      console.log(`✅ ${file} JSON syntax is valid`);
    } catch (err) {
      console.error(`❌ ${file} has invalid JSON syntax:`, err.message);
      allValid = false;
    }
  }
}

// Validate budget examples (JSON syntax only - budget schema validation separate)
const budgetsDir = "examples/budgets";
if (fs.existsSync(budgetsDir)) {
  const jsonFiles = fs
    .readdirSync(budgetsDir)
    .filter((f) => f.endsWith(".json"));

  for (const file of jsonFiles) {
    const filePath = path.join(budgetsDir, file);
    try {
      JSON.parse(fs.readFileSync(filePath, "utf8"));
      console.log(`✅ ${file} budget JSON syntax is valid`);
    } catch (err) {
      console.error(`❌ ${file} budget has invalid JSON syntax:`, err.message);
      allValid = false;
    }
  }
}

// Load and compile budget request schema
const budgetSchemaPath = "schemas/budget_spec.schema.json";
let validateBudgetRequest;
if (fs.existsSync(budgetSchemaPath)) {
  try {
    const budgetSchema = JSON.parse(fs.readFileSync(budgetSchemaPath, "utf8"));
    const cleanBudgetSchema = { ...budgetSchema };
    delete cleanBudgetSchema.$schema;
    validateBudgetRequest = ajv.compile(cleanBudgetSchema);
    console.log("✅ Budget request schema compilation successful");
  } catch (err) {
    console.error("❌ Budget request schema compilation error:", err.message);
    allValid = false;
  }
}

// Validate budget request examples against schema
const budgetRequestsDir = "examples/requests";
if (fs.existsSync(budgetRequestsDir) && validateBudgetRequest) {
  const jsonFiles = fs
    .readdirSync(budgetRequestsDir)
    .filter((f) => f.includes("budget") && f.endsWith(".json"));

  for (const file of jsonFiles) {
    const filePath = path.join(budgetRequestsDir, file);
    try {
      const data = JSON.parse(fs.readFileSync(filePath, "utf8"));
      const valid = validateBudgetRequest(data);

      if (valid) {
        console.log(`✅ ${file} request is valid`);
      } else {
        console.error(`❌ ${file} request is invalid:`);
        console.error(ajv.errorsText(validateBudgetRequest.errors));
        allValid = false;
      }
    } catch (err) {
      console.error(`❌ ${file} request has invalid JSON syntax:`, err.message);
      allValid = false;
    }
  }
}

if (!allValid) {
  process.exit(1);
}

console.log("🎉 All examples are valid!");
