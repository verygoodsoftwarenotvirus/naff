# Meal Planning System Overview

## Core Concepts

The meal planning system is designed to help groups collaboratively decide what to eat through a democratic voting process. Here's how the key concepts relate:

- **[Recipe](recipes.md)**: A detailed cooking instruction with ingredients, steps, tools, and timing
- **[Meal](meals.md)**: A collection of one or more recipes that form a complete eating experience (e.g., "Dinner" = main course + side dish)
- **Meal Plan Event**: A specific time slot for eating (e.g., "Dinner on Tuesday")
- **Meal Plan Option**: A proposed meal for a specific event (e.g., "Pasta with salad for Tuesday dinner")
- **Meal Plan**: A collection of events with their options, covering a time period (e.g., "This week's meal plan")

## System Architecture

The meal planning system consists of several key components:

### 1. Core Domain Models

- **Meal Planning Manager**: Handles CRUD operations for meals, meal plans, and related entities
- **Recipe Manager**: Manages recipe creation, updates, and retrieval
- **Valid Enumerations Manager**: Manages system-defined valid values (ingredients, measurement units, etc.)

### 2. Voting and Decision Making

- **Schulze Voting Method**: Used for ranking meal plan options within each event
- **Vote Management**: Users can rank options and change votes until the deadline
- **Election Processing**: Background workers tally votes and determine winners

### 3. Background Workers

- **Meal Plan Finalizer**: Processes finalized meal plans and creates grocery lists
- **Grocery List Initializer**: Generates shopping lists from finalized meal plans
- **Task Creator**: Creates prep tasks and cooking assignments
- **Search Data Index Scheduler**: Updates search indices for recipe discovery

## Meal Planning Flow

### 1. Meal Plan Creation

1. User creates a meal plan with a name and time period (ad-hoc, at will)
2. Events are added to the meal plan (e.g., "Dinner on Tuesday")
3. For each event, multiple meal options are proposed
4. Each option references a meal (collection of recipes) - creating users can search for existing meals to add

### 2. Voting Phase

1. Users rank their preferences for each event's options
2. Users can change their votes until the deadline
3. Users can abstain from voting on some options while voting on others
4. The system tracks who has and hasn't voted
5. If no one votes on any options, all options are considered tied and one is chosen and marked as tiebroken

### 3. Finalization

1. When the voting deadline passes, the system automatically finalizes the meal plan
2. The Schulze voting method determines the winning option for each event
3. Background workers process the finalized meal plan:
   - Create grocery lists from all winning recipes
   - Generate prep tasks based on recipe requirements
   - Assign cooking responsibilities

### 4. Execution

1. Users can view their assigned tasks and grocery lists
2. Grocery lists can be modified by any account member (mark items as acquired, edit quantities, etc.)
3. Prep tasks can be completed ahead of time
4. Cooking tasks are completed on the scheduled day

## Key Design Decisions

### Why Schulze Voting?

The Schulze method was chosen after research comparing various election tallying schemes. While no voting system can be perfectly fair in all circumstances, Schulze was found to be the closest to ideal. It handles complex preference rankings well and produces consistent, defensible results.

**TODO**: Implement instant runoff voting as a fallback for very small groups (2-3 people) where Schulze doesn't work as effectively.

### Why Background Workers?

The system uses background workers to handle computationally intensive tasks like:

- Vote tallying and election processing
- Grocery list generation and ingredient consolidation
- Task creation and assignment
- Search index updates

This keeps the API responsive while ensuring complex operations complete reliably.

## Data Models

### Meal Plan Events

Events represent specific eating times within a meal plan. They include:

- **Meal Name**: Enum values like "breakfast", "lunch", "dinner", "snack"
- **Start/End Times**: Flexible timing to accommodate busy family schedules
- **Notes**: Additional context (e.g., "Sarah has volleyball, so dinner is at 7:30")

### Meal Plan Options

Options are proposed meals for specific events:

- **Meal Reference**: Points to a [meal](meals.md) (collection of recipes)
- **Assigned Cook/Dishwasher**: Optional role assignments
- **Notes**: Additional context (e.g. "this was a hit last time")

### Voting System

- **Ranking**: Users rank options in order of preference
- **Abstention**: Users can abstain from voting on some options
- **Vote Changes**: Votes can be modified until the deadline
- **Deadline Enforcement**: Votes can be changed after finalization, but changes have no effect on the already-determined winners

**TODO**: Add validation to prevent vote changes after meal plan finalization to avoid user confusion.

## Recipe Option Selections

Recipes can include [option groups](recipes.md#option-groups-alternative-ingredients-instruments-and-vessels) - alternative ingredients, instruments, or vessels where any one can be used. The meal planning system allows users to specify their preferences for these alternatives through **selections**.

### Selection Fields

When creating a meal plan, users can include selections that specify which alternative to use:

- **Selection Type**: Can be `ingredient`, `instrument`, or `vessel`
- **Recipe ID**: The recipe containing the option group
- **Recipe Step ID**: The step within the recipe containing the option group
- **Ingredient Index**: Which option group in the step (maps to the `Index` field in the recipe)
- **Selected Option Index**: The user's chosen alternative (maps to the `OptionIndex` field in the recipe)

### Grocery List Generation Behavior

When generating grocery lists from finalized meal plans, the system handles selections as follows:

1. **With User Selection**: If the user specified a selection for an option group, only the selected alternative appears in the grocery list
2. **Without Selection (Default)**: If no selection was made, the system defaults to `optionIndex: 0` (the first alternative in the group)

This ensures grocery lists are practical and actionable - users won't see both "butter" and "margarine" for the same recipe step.

### Creating Selections at Meal Plan Creation

Selections can be provided when creating a meal plan:

```json
{
  "votingDeadline": "2024-01-15T18:00:00Z",
  "events": [...],
  "selections": [
    {
      "recipeId": "abc123",
      "recipeStepId": "step456",
      "ingredientIndex": 0,
      "selectedOptionIndex": 1,
      "selectionType": "ingredient"
    }
  ]
}
```

This example selects the second alternative (`optionIndex: 1`, e.g., margarine) for the first option group (`ingredientIndex: 0`) in the specified recipe step.

### Managing Selections

Selections can be:

- **Created**: When initially specifying preferences
- **Updated**: To change the selected option index
- **Archived**: To remove a selection (reverts to default behavior)

### Design Rationale

The selection system was designed with these principles:

1. **Sensible Defaults**: Without explicit selections, the first option (`optionIndex: 0`) is used - typically the "primary" or most common ingredient
2. **User Control**: Users who care about specific alternatives can override defaults
3. **Clean Grocery Lists**: Only selected alternatives appear, avoiding confusion
4. **Flexibility**: Supports ingredients, instruments, and vessels to cover all recipe variation scenarios

## Integration Points

### User Ingredient Preferences

Users can set preferences for ingredients they prefer or prefer not to eat. This data is currently not used by the system but is intended for future features like:

- Warning users about recipes containing disliked ingredients
- Highlighting recipes featuring preferred ingredients

**TODO**: Implement recipe filtering based on user ingredient preferences.

### Recipe Analysis

The system analyzes [recipes](recipes.md) to:

- Identify prep tasks that can be done ahead of time
- Determine cooking dependencies and timing
- Generate appropriate task assignments

## Current Limitations

1. **Cron-based Processing**: All workers run on schedules, even when there's no work
2. **No Real-time Updates**: Users must refresh to see finalization results
3. **Limited Election Methods**: Only Schulze is implemented (Instant Runoff is defined but not used)
4. **Manual Task Assignment**: Prep tasks aren't automatically assigned to users
5. **No UI**: This is a backend API service only - access via gRPC/HTTP endpoints (e.g., Postman)
6. **Recipe Management**: Only service admins can create recipes; users can clone existing recipes

## Future Improvements

### Known Edge Cases

- **Recipe Deletion**: If a [recipe](recipes.md) is deleted after a meal plan is created but before finalization, it could permanently break the meal plan
- **Account Membership Changes**: If a user leaves an account after voting but before finalization, their vote is still counted
- **Overlapping Events**: The system doesn't prevent overlapping meal times

**TODO**: Add validation to prevent recipe deletion when referenced by active meal plans.
**TODO**: Add validation to prevent meal plan modifications after finalization.
**TODO**: Add integration tests for account membership changes during voting.

### Planned Improvements

1. **Queue-Based Architecture**: Move from cron jobs to queue-based processing for better scalability
2. **Enhanced Validation**: Add comprehensive checks for recipe/meal/account interactions
3. **Grocery List Consolidation**: Build interfaces for combining similar ingredients
4. **Notification System**: Ensure user notifications are sent when meal plans are created and finalized

**TODO**: Add integration tests to verify notification behavior.
**TODO**: Add integration tests to verify email sending behavior.

## Testing

The system includes comprehensive integration tests covering:

- Complete meal plan lifecycle (creation, voting, finalization)
- [Recipe](recipes.md) management and validation
- User ingredient preferences
- Valid enumeration management
- Voting logic and election processing

**TODO**: Add performance tests for background workers.
**TODO**: Add integration tests for edge cases around account membership and recipe deletion.
