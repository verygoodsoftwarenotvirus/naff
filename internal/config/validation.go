package config

import (
	"fmt"
	"strings"
)

// validFieldTypes is the set of supported Go types for entity fields.
var validFieldTypes = map[string]bool{
	"bool":    true,
	"string":  true,
	"int":     true,
	"int8":    true,
	"int16":   true,
	"int32":   true,
	"int64":   true,
	"uint":    true,
	"uint8":   true,
	"uint16":  true,
	"uint32":  true,
	"uint64":  true,
	"float32": true,
	"float64": true,
	// time.Time and *time.Time map to TIMESTAMPTZ columns; the platform's
	// Null*/*Pointer* helpers cover the nullable conversions.
	"time.Time": true,
}

// Validate checks the project config for errors.
func (p *Project) Validate() error {
	if p.ProjectMeta.Name == "" {
		return fmt.Errorf("project.name is required")
	}

	if p.Targets.BackendEnabled() && p.ProjectMeta.Module == "" {
		return fmt.Errorf("project.module is required when backend target is enabled")
	}

	if p.Targets.IOS {
		// Default IOSModuleName to project name if not set.
		if p.ProjectMeta.IOSModuleName == "" {
			p.ProjectMeta.IOSModuleName = p.ProjectMeta.Name
		}
		// Default IOSBundleID if not set.
		if p.ProjectMeta.IOSBundleID == "" {
			p.ProjectMeta.IOSBundleID = "com.example." + strings.ToLower(p.ProjectMeta.Name)
		}
	}

	if len(p.Domains) == 0 {
		return fmt.Errorf("at least one domain is required")
	}

	// Build a set of all entity names for ownership validation.
	allEntities := make(map[string]string) // entity name -> domain name
	for _, d := range p.Domains {
		if d.Name == "" {
			return fmt.Errorf("domain name is required")
		}
		if len(d.Entities) == 0 {
			return fmt.Errorf("domain %q must have at least one entity", d.Name)
		}
		for _, e := range d.Entities {
			if e.Name == "" {
				return fmt.Errorf("entity name is required in domain %q", d.Name)
			}
			if _, exists := allEntities[e.Name]; exists {
				return fmt.Errorf("duplicate entity name %q", e.Name)
			}
			allEntities[e.Name] = d.Name
		}
	}

	// Validate each entity.
	for _, d := range p.Domains {
		for _, e := range d.Entities {
			if err := validateEntity(e, allEntities); err != nil {
				return fmt.Errorf("domain %q, entity %q: %w", d.Name, e.Name, err)
			}
		}
	}

	// Check for ownership cycles.
	if err := checkOwnershipCycles(p.Domains); err != nil {
		return err
	}

	return nil
}

func validateEntity(e Entity, allEntities map[string]string) error {
	if len(e.Fields) == 0 {
		return fmt.Errorf("must have at least one field")
	}

	if e.BelongsTo != "" {
		if _, ok := allEntities[e.BelongsTo]; !ok {
			return fmt.Errorf("belongs_to references unknown entity %q", e.BelongsTo)
		}
	}

	for _, f := range e.Fields {
		if f.Name == "" {
			return fmt.Errorf("field name is required")
		}
		if f.Type == "" {
			return fmt.Errorf("field %q: type is required", f.Name)
		}

		baseType := f.BaseType()
		if !validFieldTypes[baseType] {
			return fmt.Errorf("field %q: unsupported type %q", f.Name, f.Type)
		}
	}

	return nil
}

func checkOwnershipCycles(domains []Domain) error {
	// Build adjacency: child -> parent
	parentOf := make(map[string]string)
	for _, d := range domains {
		for _, e := range d.Entities {
			if e.BelongsTo != "" {
				parentOf[e.Name] = e.BelongsTo
			}
		}
	}

	// Walk each chain, detect cycles.
	for entity := range parentOf {
		visited := make(map[string]bool)
		current := entity
		for {
			if visited[current] {
				// Build the cycle path for a useful error message.
				var path []string
				walk := entity
				for {
					path = append(path, walk)
					if walk == current && len(path) > 1 {
						break
					}
					next, ok := parentOf[walk]
					if !ok {
						break
					}
					walk = next
				}
				return fmt.Errorf("ownership cycle detected: %s", strings.Join(path, " -> "))
			}
			visited[current] = true
			parent, ok := parentOf[current]
			if !ok {
				break
			}
			current = parent
		}
	}

	return nil
}
