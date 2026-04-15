//
//  LocalDataProvider.swift
//  ios
//
//  Provides offline data from a bundled JSON seed file.
//  Decodes Go-serialized JSON into Codable types, then converts to Protobuf types used by the UI.
//

import Foundation
import SwiftProtobuf

// MARK: - LocalDataProvider

@MainActor
class LocalDataProvider {
  static let shared = LocalDataProvider()

  private(set) var recipes: [Mealplanning_Recipe] = []
  private(set) var meals: [Mealplanning_Meal] = []
  private(set) var validIngredients: [SeedValidIngredient] = []
  private(set) var validPreparations: [SeedValidPreparation] = []

  private var isLoaded = false

  func loadIfNeeded() {
    guard !isLoaded else { return }
    isLoaded = true

    guard let url = Bundle.main.url(forResource: "seed_data", withExtension: "json") else {
      print("seed_data.json not found in bundle")
      return
    }

    do {
      let data = try Data(contentsOf: url)
      let decoder = JSONDecoder()
      decoder.dateDecodingStrategy = .custom { decoder in
        let container = try decoder.singleValueContainer()
        let dateString = try container.decode(String.self)
        if let date = ISO8601DateFormatter().date(from: dateString) {
          return date
        }
        // Try with fractional seconds
        let formatter = ISO8601DateFormatter()
        formatter.formatOptions.insert(.withFractionalSeconds)
        if let date = formatter.date(from: dateString) {
          return date
        }
        throw DecodingError.dataCorruptedError(
          in: container, debugDescription: "Cannot decode date: \(dateString)")
      }

      let seed = try decoder.decode(SeedData.self, from: data)

      // Convert to protobuf types
      self.recipes = seed.recipes.map { $0.toProtobuf() }
      self.meals = seed.meals.map { $0.toProtobuf() }
      self.validIngredients = seed.enumerations.validIngredients
      self.validPreparations = seed.enumerations.validPreparations

      print("LocalDataProvider: loaded \(recipes.count) recipes, \(meals.count) meals")
    } catch {
      print("LocalDataProvider: failed to load seed data: \(error)")
    }
  }

  // MARK: - Query Methods

  func getRecipes() -> [Mealplanning_Recipe] {
    return recipes
  }

  func getRecipe(id: String) -> Mealplanning_Recipe? {
    return recipes.first { $0.id == id }
  }

  func searchRecipes(query: String) -> [Mealplanning_Recipe] {
    let q = query.lowercased()
    return recipes.filter {
      $0.name.lowercased().contains(q) || $0.description_p.lowercased().contains(q)
    }
  }

  func getMeals() -> [Mealplanning_Meal] {
    return meals
  }

  func getMeal(id: String) -> Mealplanning_Meal? {
    return meals.first { $0.id == id }
  }

  func searchMeals(query: String) -> [Mealplanning_Meal] {
    let q = query.lowercased()
    return meals.filter {
      $0.name.lowercased().contains(q) || $0.description_p.lowercased().contains(q)
    }
  }
}

// MARK: - Seed Data Codable Types

struct SeedData: Codable {
  let exportedAt: Date
  let enumerations: SeedEnumerations
  let recipes: [SeedRecipe]
  let meals: [SeedMeal]
}

struct SeedEnumerations: Codable {
  let validIngredients: [SeedValidIngredient]
  let validPreparations: [SeedValidPreparation]
  let validInstruments: [SeedValidInstrument]
  let validVessels: [SeedValidVessel]
  let validMeasurementUnits: [SeedValidMeasurementUnit]
  let validIngredientStates: [SeedValidIngredientState]
  let validIngredientPreparations: [SeedValidIngredientPreparation]
  let validIngredientMeasurementUnits: [SeedValidIngredientMeasurementUnit]
  let validPreparationInstruments: [SeedValidPreparationInstrument]
  let validPreparationVessels: [SeedValidPreparationVessel]
  let validIngredientGroups: [SeedValidIngredientGroup]
  let validIngredientStateIngredients: [SeedValidIngredientStateIngredient]
  let validMeasurementUnitConversions: [SeedValidMeasurementUnitConversion]
}

// MARK: - Base enumeration types

struct SeedValidIngredient: Codable {
  let id: String
  let name: String
  let description: String
  let slug: String
  let pluralName: String
  let warning: String?
  let iconPath: String?
  let containsEgg: Bool
  let containsDairy: Bool
  let containsPeanut: Bool
  let containsTreeNut: Bool
  let containsSoy: Bool
  let containsWheat: Bool
  let containsShellfish: Bool
  let containsSesame: Bool
  let containsFish: Bool
  let containsGluten: Bool
  let containsAlcohol: Bool
  let animalFlesh: Bool
  let animalDerived: Bool
  let isLiquid: Bool
  let isStarch: Bool
  let isProtein: Bool
  let isGrain: Bool
  let isFruit: Bool
  let isSalt: Bool
  let isFat: Bool
  let isAcid: Bool
  let isHeat: Bool
  let restrictToPreparations: Bool
  let contaminatesEquipment: Bool
  let storageInstructions: String?
  let shoppingSuggestions: String?
  let storageTemperatureInCelsius: SeedOptionalFloat32Range?

  func toProtobuf() -> Mealplanning_ValidIngredient {
    var pb = Mealplanning_ValidIngredient()
    pb.id = id
    pb.name = name
    pb.description_p = description
    pb.slug = slug
    pb.pluralName = pluralName
    pb.warning = warning ?? ""
    pb.iconPath = iconPath ?? ""
    pb.containsEgg = containsEgg
    pb.containsDairy = containsDairy
    pb.containsPeanut = containsPeanut
    pb.containsTreeNut = containsTreeNut
    pb.containsSoy = containsSoy
    pb.containsWheat = containsWheat
    pb.containsShellfish = containsShellfish
    pb.containsSesame = containsSesame
    pb.containsFish = containsFish
    pb.containsGluten = containsGluten
    pb.containsAlcohol = containsAlcohol
    pb.animalFlesh = animalFlesh
    pb.animalDerived = animalDerived
    pb.isLiquid = isLiquid
    pb.isStarch = isStarch
    pb.isProtein = isProtein
    pb.isGrain = isGrain
    pb.isFruit = isFruit
    pb.isSalt = isSalt
    pb.isFat = isFat
    pb.isAcid = isAcid
    pb.isHeat = isHeat
    pb.restrictToPreparations = restrictToPreparations
    pb.contaminatesEquipment = contaminatesEquipment
    pb.storageInstructions = storageInstructions ?? ""
    pb.shoppingSuggestions = shoppingSuggestions ?? ""
    return pb
  }
}

struct SeedValidPreparation: Codable {
  let id: String
  let name: String
  let description: String
  let slug: String
  let pastTense: String
  let iconPath: String?
  let temperatureRequired: Bool
  let timeEstimateRequired: Bool
  let conditionExpressionRequired: Bool
  let consumesVessel: Bool
  let onlyForVessels: Bool
  let restrictToIngredients: Bool
  let yieldsNothing: Bool

  func toProtobuf() -> Mealplanning_ValidPreparation {
    var pb = Mealplanning_ValidPreparation()
    pb.id = id
    pb.name = name
    pb.description_p = description
    pb.slug = slug
    pb.pastTense = pastTense
    pb.iconPath = iconPath ?? ""
    pb.temperatureRequired = temperatureRequired
    pb.timeEstimateRequired = timeEstimateRequired
    pb.conditionExpressionRequired = conditionExpressionRequired
    pb.consumesVessel = consumesVessel
    pb.onlyForVessels = onlyForVessels
    pb.restrictToIngredients = restrictToIngredients
    pb.yieldsNothing = yieldsNothing
    return pb
  }
}

struct SeedValidInstrument: Codable {
  let id: String
  let name: String
  let pluralName: String
  let description: String
  let slug: String
  let iconPath: String?
  let displayInSummaryLists: Bool
  let usableForStorage: Bool
  let includeInGeneratedInstructions: Bool

  func toProtobuf() -> Mealplanning_ValidInstrument {
    var pb = Mealplanning_ValidInstrument()
    pb.id = id
    pb.name = name
    pb.pluralName = pluralName
    pb.description_p = description
    pb.slug = slug
    pb.iconPath = iconPath ?? ""
    pb.displayInSummaryLists = displayInSummaryLists
    pb.usableForStorage = usableForStorage
    pb.includeInGeneratedInstructions = includeInGeneratedInstructions
    return pb
  }
}

struct SeedValidVessel: Codable {
  let id: String
  let name: String
  let pluralName: String
  let description: String
  let slug: String
  let iconPath: String?
  let shape: String
  let capacity: Float
  let capacityUnit: SeedValidMeasurementUnit?
  let widthInMillimeters: Float
  let lengthInMillimeters: Float
  let heightInMillimeters: Float
  let displayInSummaryLists: Bool
  let usableForStorage: Bool
  let includeInGeneratedInstructions: Bool

  func toProtobuf() -> Mealplanning_ValidVessel {
    var pb = Mealplanning_ValidVessel()
    pb.id = id
    pb.name = name
    pb.pluralName = pluralName
    pb.description_p = description
    pb.slug = slug
    pb.iconPath = iconPath ?? ""
    pb.shape = validVesselShapeFromString(shape)
    pb.capacity = capacity
    if let cu = capacityUnit { pb.capacityUnit = cu.toProtobuf() }
    pb.widthInMillimeters = widthInMillimeters
    pb.lengthInMillimeters = lengthInMillimeters
    pb.heightInMillimeters = heightInMillimeters
    pb.displayInSummaryLists = displayInSummaryLists
    pb.usableForStorage = usableForStorage
    pb.includeInGeneratedInstructions = includeInGeneratedInstructions
    return pb
  }
}

struct SeedValidMeasurementUnit: Codable {
  let id: String
  let name: String
  let description: String
  let slug: String
  let pluralName: String
  let iconPath: String?
  let volumetric: Bool
  let universal: Bool
  let metric: Bool
  let imperial: Bool

  func toProtobuf() -> Mealplanning_ValidMeasurementUnit {
    var pb = Mealplanning_ValidMeasurementUnit()
    pb.id = id
    pb.name = name
    pb.description_p = description
    pb.slug = slug
    pb.pluralName = pluralName
    pb.iconPath = iconPath ?? ""
    pb.volumetric = volumetric
    pb.universal = universal
    pb.metric = metric
    pb.imperial = imperial
    return pb
  }
}

struct SeedValidIngredientState: Codable {
  let id: String
  let name: String
  let slug: String
  let pastTense: String
  let description: String
  let attributeType: String
  let iconPath: String?

  func toProtobuf() -> Mealplanning_ValidIngredientState {
    var pb = Mealplanning_ValidIngredientState()
    pb.id = id
    pb.name = name
    pb.slug = slug
    pb.pastTense = pastTense
    pb.description_p = description
    pb.attributeType = ingredientStateAttributeTypeFromString(attributeType)
    pb.iconPath = iconPath ?? ""
    return pb
  }
}

// MARK: - Bridge types (not converted to protobuf since UI doesn't use them directly)

struct SeedValidIngredientPreparation: Codable {
  let id: String
  let notes: String
  let preparation: SeedValidPreparation
  let ingredient: SeedValidIngredient
}

struct SeedValidIngredientMeasurementUnit: Codable {
  let id: String
  let notes: String
  let allowableQuantity: SeedFloat32RangeWithOptionalMax?
  let measurementUnit: SeedValidMeasurementUnit
  let ingredient: SeedValidIngredient
}

struct SeedValidPreparationInstrument: Codable {
  let id: String
  let notes: String
  let preparation: SeedValidPreparation
  let instrument: SeedValidInstrument
}

struct SeedValidPreparationVessel: Codable {
  let id: String
  let notes: String
  let preparation: SeedValidPreparation
  let vessel: SeedValidVessel
}

struct SeedValidIngredientGroup: Codable {
  let id: String
  let name: String
  let slug: String
  let description: String
  let members: [SeedValidIngredientGroupMember]
}

struct SeedValidIngredientGroupMember: Codable {
  let id: String
  let belongsToGroup: String
  let validIngredient: SeedValidIngredient
}

struct SeedValidIngredientStateIngredient: Codable {
  let id: String
  let notes: String
  let ingredientState: SeedValidIngredientState
  let ingredient: SeedValidIngredient
}

struct SeedValidMeasurementUnitConversion: Codable {
  let id: String
  let from: SeedValidMeasurementUnit
  let to: SeedValidMeasurementUnit
  let notes: String
  let modifier: Float
  let onlyForIngredient: SeedValidIngredient?
}

// MARK: - Range types

struct SeedOptionalFloat32Range: Codable {
  let min: Float?
  let max: Float?
}

struct SeedFloat32RangeWithOptionalMax: Codable {
  let min: Float
  let max: Float?

  func toProtobuf() -> Common_Float32RangeWithOptionalMax {
    var pb = Common_Float32RangeWithOptionalMax()
    pb.min = min
    if let m = max { pb.max = m }
    return pb
  }
}

struct SeedOptionalUint32Range: Codable {
  let min: UInt32?
  let max: UInt32?

  func toProtobuf() -> Common_OptionalUint32Range {
    var pb = Common_OptionalUint32Range()
    if let m = min { pb.min = m }
    if let m = max { pb.max = m }
    return pb
  }
}

struct SeedOptionalFloat32RangePB: Codable {
  let min: Float?
  let max: Float?

  func toProtobuf() -> Common_OptionalFloat32Range {
    var pb = Common_OptionalFloat32Range()
    if let m = min { pb.min = m }
    if let m = max { pb.max = m }
    return pb
  }
}

struct SeedUint32RangeWithOptionalMax: Codable {
  let min: UInt32
  let max: UInt32?

  func toProtobuf() -> Common_Uint32RangeWithOptionalMax {
    var pb = Common_Uint32RangeWithOptionalMax()
    pb.min = min
    if let m = max { pb.max = m }
    return pb
  }
}

struct SeedUint16RangeWithOptionalMax: Codable {
  let min: UInt16
  let max: UInt16?
}

// MARK: - Recipe types

struct SeedRecipe: Codable {
  let id: String
  let name: String
  let slug: String
  let source: String
  let sourceISBN: String?
  let description: String
  let portionName: String
  let pluralPortionName: String
  let createdByUser: String
  let yieldsComponentType: String?
  let status: String
  let inspiredByRecipeID: String?
  let estimatedPortions: SeedFloat32RangeWithOptionalMax
  let eligibleForMeals: Bool
  let steps: [SeedRecipeStep]
  let media: [SeedRecipeMedia]?
  let prepTasks: [SeedRecipePrepTask]?
  let associatedRecipes: [SeedRecipe]?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_Recipe {
    var pb = Mealplanning_Recipe()
    pb.id = id
    pb.name = name
    pb.slug = slug
    pb.source = source
    pb.description_p = description
    pb.portionName = portionName
    pb.pluralPortionName = pluralPortionName
    pb.createdByUser = createdByUser
    pb.yieldsComponentType = mealComponentTypeFromString(yieldsComponentType ?? "")
    pb.status = status
    if let rid = inspiredByRecipeID { pb.inspiredByRecipeID = rid }
    pb.estimatedPortions = estimatedPortions.toProtobuf()
    pb.eligibleForMeals = eligibleForMeals
    pb.steps = steps.map { $0.toProtobuf() }
    pb.media = (media ?? []).map { $0.toProtobuf() }
    pb.prepTasks = (prepTasks ?? []).map { $0.toProtobuf() }
    pb.associatedRecipes = (associatedRecipes ?? []).map { $0.toProtobuf() }
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStep: Codable {
  let id: String
  let index: UInt32
  let belongsToRecipe: String
  let preparation: SeedValidPreparation
  let notes: String
  let explicitInstructions: String?
  let conditionExpression: String?
  let estimatedTimeInSeconds: SeedOptionalUint32Range?
  let temperatureInCelsius: SeedOptionalFloat32Range?
  let optional: Bool
  let startTimerAutomatically: Bool
  let ingredients: [SeedRecipeStepIngredient]
  let instruments: [SeedRecipeStepInstrument]
  let vessels: [SeedRecipeStepVessel]
  let products: [SeedRecipeStepProduct]
  let completionConditions: [SeedRecipeStepCompletionCondition]?
  let media: [SeedRecipeMedia]?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStep {
    var pb = Mealplanning_RecipeStep()
    pb.id = id
    pb.index = index
    pb.belongsToRecipe = belongsToRecipe
    pb.preparation = preparation.toProtobuf()
    pb.notes = notes
    pb.explicitInstructions = explicitInstructions ?? ""
    pb.conditionExpression = conditionExpression ?? ""
    if let t = estimatedTimeInSeconds { pb.estimatedTimeInSeconds = t.toProtobuf() }
    if let t = temperatureInCelsius {
      var r = Common_OptionalFloat32Range()
      if let m = t.min { r.min = m }
      if let m = t.max { r.max = m }
      pb.temperatureInCelsius = r
    }
    pb.optional = `optional`
    pb.startTimerAutomatically = startTimerAutomatically
    pb.ingredients = ingredients.map { $0.toProtobuf() }
    pb.instruments = instruments.map { $0.toProtobuf() }
    pb.vessels = vessels.map { $0.toProtobuf() }
    pb.products = products.map { $0.toProtobuf() }
    pb.completionConditions = (completionConditions ?? []).map { $0.toProtobuf() }
    pb.media = (media ?? []).map { $0.toProtobuf() }
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStepIngredient: Codable {
  let id: String
  let belongsToRecipeStep: String
  let name: String
  let ingredient: SeedValidIngredient?
  let measurementUnit: SeedValidMeasurementUnit
  let quantity: SeedFloat32RangeWithOptionalMax
  let quantityNotes: String?
  let ingredientNotes: String?
  let productOfRecipeID: String?
  let recipeStepProductID: String?
  let vesselIndex: UInt16?
  let productPercentageToUse: Float?
  let index: UInt32
  let optionIndex: UInt32?
  let optional: Bool
  let toTaste: Bool
  let scaleFactor: Float?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepIngredient {
    var pb = Mealplanning_RecipeStepIngredient()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.name = name
    if let ing = ingredient { pb.ingredient = ing.toProtobuf() }
    pb.measurementUnit = measurementUnit.toProtobuf()
    pb.quantity = quantity.toProtobuf()
    pb.quantityNotes = quantityNotes ?? ""
    pb.ingredientNotes = ingredientNotes ?? ""
    if let rid = productOfRecipeID { pb.recipeStepProductRecipeID = rid }
    if let pid = recipeStepProductID { pb.recipeStepProductID = pid }
    if let vi = vesselIndex { pb.vesselIndex = UInt32(vi) }
    if let ppu = productPercentageToUse { pb.productPercentageToUse = ppu }
    pb.index = index
    pb.optionIndex = optionIndex ?? 0
    pb.optional = `optional`
    pb.toTaste = toTaste
    pb.scaleFactor = scaleFactor ?? 1.0
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStepInstrument: Codable {
  let id: String
  let belongsToRecipeStep: String
  let name: String
  let instrument: SeedValidInstrument?
  let notes: String?
  let recipeStepProductID: String?
  let quantity: SeedUint32RangeWithOptionalMax?
  let index: UInt32
  let optionIndex: UInt32?
  let preferenceRank: UInt32?
  let optional: Bool
  let scaleFactor: Float?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepInstrument {
    var pb = Mealplanning_RecipeStepInstrument()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.name = name
    if let inst = instrument { pb.instrument = inst.toProtobuf() }
    pb.notes = notes ?? ""
    if let pid = recipeStepProductID { pb.recipeStepProductID = pid }
    if let q = quantity { pb.quantity = q.toProtobuf() }
    pb.index = index
    pb.optionIndex = optionIndex ?? 0
    pb.preferenceRank = preferenceRank ?? 0
    pb.optional = `optional`
    pb.scaleFactor = scaleFactor ?? 1.0
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStepVessel: Codable {
  let id: String
  let belongsToRecipeStep: String
  let name: String
  let vessel: SeedValidVessel?
  let vesselPreposition: String?
  let notes: String?
  let recipeStepProductID: String?
  let quantity: SeedUint16RangeWithOptionalMax?
  let index: UInt32
  let optionIndex: UInt32?
  let unavailableAfterStep: Bool
  let scaleFactor: Float?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepVessel {
    var pb = Mealplanning_RecipeStepVessel()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.name = name
    if let v = vessel { pb.vessel = v.toProtobuf() }
    pb.vesselPreposition = vesselPreposition ?? ""
    pb.notes = notes ?? ""
    if let pid = recipeStepProductID { pb.recipeStepProductID = pid }
    if let q = quantity {
      var r = Common_Uint16RangeWithOptionalMax()
      r.min = UInt32(q.min)
      if let m = q.max { r.max = UInt32(m) }
      pb.quantity = r
    }
    pb.index = index
    pb.optionIndex = optionIndex ?? 0
    pb.unavailableAfterStep = unavailableAfterStep
    pb.scaleFactor = scaleFactor ?? 1.0
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStepProduct: Codable {
  let id: String
  let belongsToRecipeStep: String
  let name: String
  let type: String
  let storageInstructions: String?
  let quantityNotes: String?
  let measurementUnit: SeedValidMeasurementUnit?
  let measurementQuantity: SeedOptionalFloat32Range?
  let itemQuantity: SeedOptionalFloat32Range?
  let containedInVesselIndex: UInt16?
  let storageTemperatureInCelsius: SeedOptionalFloat32Range?
  let storageDurationInSeconds: SeedOptionalUint32Range?
  let index: UInt32
  let isWaste: Bool
  let isLiquid: Bool
  let compostable: Bool
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepProduct {
    var pb = Mealplanning_RecipeStepProduct()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.name = name
    pb.type = recipeStepProductTypeFromString(type)
    pb.storageInstructions = storageInstructions ?? ""
    pb.quantityNotes = quantityNotes ?? ""
    if let mu = measurementUnit { pb.measurementUnit = mu.toProtobuf() }
    if let mq = measurementQuantity {
      var r = Common_OptionalFloat32Range()
      if let m = mq.min { r.min = m }
      if let m = mq.max { r.max = m }
      pb.measurementQuantity = r
    }
    if let iq = itemQuantity {
      var r = Common_OptionalFloat32Range()
      if let m = iq.min { r.min = m }
      if let m = iq.max { r.max = m }
      pb.itemQuantity = r
    }
    if let ci = containedInVesselIndex { pb.containedInVesselIndex = UInt32(ci) }
    pb.index = index
    pb.isWaste = isWaste
    pb.isLiquid = isLiquid
    pb.compostable = compostable
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeStepCompletionCondition: Codable {
  let id: String
  let belongsToRecipeStep: String
  let ingredientState: SeedValidIngredientState
  let notes: String?
  let ingredients: [SeedStepConditionIngredient]
  let optional: Bool
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepCompletionCondition {
    var pb = Mealplanning_RecipeStepCompletionCondition()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.ingredientState = ingredientState.toProtobuf()
    pb.notes = notes ?? ""
    pb.ingredients = ingredients.map { $0.toProtobuf() }
    pb.optional = `optional`
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedStepConditionIngredient: Codable {
  let id: String
  let belongsToRecipeStepCompletionCondition: String
  let recipeStepIngredient: String
  let createdAt: Date
  let archivedAt: Date?
  let lastUpdatedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeStepCompletionConditionIngredient {
    var pb = Mealplanning_RecipeStepCompletionConditionIngredient()
    pb.id = id
    pb.belongsToRecipeStepCompletionCondition = belongsToRecipeStepCompletionCondition
    pb.recipeStepIngredient = recipeStepIngredient
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipeMedia: Codable {
  let id: String
  let belongsToRecipe: String?
  let belongsToRecipeStep: String?
  let mimeType: String
  let internalPath: String
  let externalPath: String
  let index: UInt32
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipeMedia {
    var pb = Mealplanning_RecipeMedia()
    pb.id = id
    if let r = belongsToRecipe { pb.belongsToRecipe = r }
    if let r = belongsToRecipeStep { pb.belongsToRecipeStep = r }
    pb.mimeType = mimeType
    pb.internalPath = internalPath
    pb.externalPath = externalPath
    pb.index = index
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipePrepTask: Codable {
  let id: String
  let belongsToRecipe: String
  let name: String
  let description: String
  let notes: String?
  let explicitStorageInstructions: String?
  let storageType: String
  let optional: Bool
  let storageTemperatureInCelsius: SeedOptionalFloat32Range?
  let timeBufferBeforeRecipeInSeconds: SeedUint32RangeWithOptionalMax?
  let recipeSteps: [SeedRecipePrepTaskStep]?
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_RecipePrepTask {
    var pb = Mealplanning_RecipePrepTask()
    pb.id = id
    pb.belongsToRecipe = belongsToRecipe
    pb.name = name
    pb.description_p = description
    pb.notes = notes ?? ""
    pb.explicitStorageInstructions = explicitStorageInstructions ?? ""
    pb.storageType = storageType
    pb.optional = `optional`
    if let t = timeBufferBeforeRecipeInSeconds {
      pb.timeBufferBeforeRecipeInSeconds = t.toProtobuf()
    }
    pb.taskSteps = (recipeSteps ?? []).map { $0.toProtobuf() }
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedRecipePrepTaskStep: Codable {
  let id: String
  let belongsToRecipeStep: String
  let belongsToRecipeStepTask: String
  let satisfiesRecipeStep: Bool

  func toProtobuf() -> Mealplanning_RecipePrepTaskStep {
    var pb = Mealplanning_RecipePrepTaskStep()
    pb.id = id
    pb.belongsToRecipeStep = belongsToRecipeStep
    pb.belongsToRecipePrepTask = belongsToRecipeStepTask
    pb.satisfiesRecipeStep = satisfiesRecipeStep
    return pb
  }
}

// MARK: - Meal types

struct SeedMeal: Codable {
  let id: String
  let name: String
  let description: String
  let createdByUser: String
  let estimatedPortions: SeedFloat32RangeWithOptionalMax
  let eligibleForMealPlans: Bool
  let components: [SeedMealComponent]
  let createdAt: Date
  let lastUpdatedAt: Date?
  let archivedAt: Date?

  func toProtobuf() -> Mealplanning_Meal {
    var pb = Mealplanning_Meal()
    pb.id = id
    pb.name = name
    pb.description_p = description
    pb.createdByUser = createdByUser
    pb.estimatedPortions = estimatedPortions.toProtobuf()
    pb.eligibleForMealPlans = eligibleForMealPlans
    pb.components = components.map { $0.toProtobuf() }
    pb.createdAt = dateToTimestamp(createdAt)
    if let d = lastUpdatedAt { pb.lastUpdatedAt = dateToTimestamp(d) }
    if let d = archivedAt { pb.archivedAt = dateToTimestamp(d) }
    return pb
  }
}

struct SeedMealComponent: Codable {
  let componentType: String
  let recipe: SeedRecipe
  let recipeScale: Float

  func toProtobuf() -> Mealplanning_MealComponent {
    var pb = Mealplanning_MealComponent()
    pb.componentType = mealComponentTypeFromString(componentType)
    pb.recipe = recipe.toProtobuf()
    pb.recipeScale = recipeScale
    return pb
  }
}

// MARK: - Helpers

private func dateToTimestamp(_ date: Date) -> SwiftProtobuf.Google_Protobuf_Timestamp {
  var ts = SwiftProtobuf.Google_Protobuf_Timestamp()
  ts.seconds = Int64(date.timeIntervalSince1970)
  ts.nanos = Int32((date.timeIntervalSince1970.truncatingRemainder(dividingBy: 1)) * 1_000_000_000)
  return ts
}

private func validVesselShapeFromString(_ s: String) -> Mealplanning_ValidVesselShape {
  switch s {
  case "hemisphere": return .vesselShapeHemisphere
  case "rectangle": return .vesselShapeRectangle
  case "cone": return .vesselShapeCone
  case "pyramid": return .vesselShapePyramid
  case "cylinder": return .vesselShapeCylinder
  case "sphere": return .vesselShapeSphere
  case "cube": return .vesselShapeCube
  case "other": return .vesselShapeOther
  default: return .vesselShapeOther
  }
}

private func ingredientStateAttributeTypeFromString(_ s: String)
  -> Mealplanning_ValidIngredientStateAttributeType
{
  switch s {
  case "texture": return .texture
  case "consistency": return .consistency
  case "temperature": return .temperature
  case "color": return .color
  case "appearance": return .appearance
  case "odor": return .odor
  case "taste": return .taste
  case "sound": return .sound
  case "other": return .other
  default: return .other
  }
}

private func mealComponentTypeFromString(_ s: String) -> Mealplanning_MealComponentType {
  switch s {
  case "amuse-bouche": return .amuseBouche
  case "appetizer": return .appetizer
  case "soup": return .soup
  case "main": return .main
  case "salad": return .salad
  case "beverage": return .beverage
  case "side": return .side
  case "dessert": return .dessert
  default: return .unspecified
  }
}

private func recipeStepProductTypeFromString(_ s: String) -> Mealplanning_RecipeStepProductType {
  switch s {
  case "ingredient": return .ingredient
  case "instrument": return .instrument
  case "vessel": return .vessel
  default: return .ingredient
  }
}
