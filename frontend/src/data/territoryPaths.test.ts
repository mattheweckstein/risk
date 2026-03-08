import { describe, it, expect } from 'vitest';
import { territoryPaths, connectionLines, continentColors, continentTerritories } from './territoryPaths';

describe('territoryPaths data integrity', () => {
  describe('territory count', () => {
    it('has all 42 territories', () => {
      expect(territoryPaths).toHaveLength(42);
    });
  });

  describe('territory IDs', () => {
    it('all territory IDs are unique', () => {
      const ids = territoryPaths.map((t) => t.id);
      const uniqueIds = new Set(ids);
      expect(uniqueIds.size).toBe(ids.length);
    });

    it('all territory IDs are non-empty strings', () => {
      territoryPaths.forEach((t) => {
        expect(t.id).toBeTruthy();
        expect(typeof t.id).toBe('string');
        expect(t.id.trim().length).toBeGreaterThan(0);
      });
    });
  });

  describe('label positions within viewBox bounds', () => {
    // ViewBox is 0-1100 wide, 0-600 tall based on the SVG setup
    it('all labelX values are within 0-1100', () => {
      territoryPaths.forEach((t) => {
        expect(t.labelX).toBeGreaterThanOrEqual(0);
        expect(t.labelX).toBeLessThanOrEqual(1100);
      });
    });

    it('all labelY values are within 0-600', () => {
      territoryPaths.forEach((t) => {
        expect(t.labelY).toBeGreaterThanOrEqual(0);
        expect(t.labelY).toBeLessThanOrEqual(600);
      });
    });
  });

  describe('SVG paths', () => {
    it('all territories have non-empty path strings', () => {
      territoryPaths.forEach((t) => {
        expect(t.path).toBeTruthy();
        expect(typeof t.path).toBe('string');
        expect(t.path.length).toBeGreaterThan(10);
      });
    });
  });

  describe('connection lines', () => {
    const allTerritoryIds = new Set(territoryPaths.map((t) => t.id));

    it('all connection line "from" references valid territory IDs', () => {
      connectionLines.forEach((line) => {
        expect(allTerritoryIds.has(line.from)).toBe(true);
      });
    });

    it('all connection line "to" references valid territory IDs', () => {
      connectionLines.forEach((line) => {
        expect(allTerritoryIds.has(line.to)).toBe(true);
      });
    });

    it('all connection lines have non-empty path strings', () => {
      connectionLines.forEach((line) => {
        expect(line.path).toBeTruthy();
        expect(typeof line.path).toBe('string');
      });
    });
  });

  describe('continent colors', () => {
    it('defines colors for all 6 continents', () => {
      const continents = Object.keys(continentColors);
      expect(continents).toHaveLength(6);
      expect(continents).toContain('north_america');
      expect(continents).toContain('south_america');
      expect(continents).toContain('europe');
      expect(continents).toContain('africa');
      expect(continents).toContain('asia');
      expect(continents).toContain('australia');
    });

    it('all continent colors are valid rgba strings', () => {
      Object.values(continentColors).forEach((color) => {
        expect(color).toMatch(/^rgba\(\d+,\s*\d+,\s*\d+,\s*[\d.]+\)$/);
      });
    });
  });

  describe('continent territory groupings', () => {
    const allTerritoryIds = new Set(territoryPaths.map((t) => t.id));

    it('defines groupings for all 6 continents', () => {
      expect(Object.keys(continentTerritories)).toHaveLength(6);
    });

    it('all territories in continent groupings reference valid territory IDs', () => {
      Object.entries(continentTerritories).forEach(([continent, territories]) => {
        territories.forEach((tid) => {
          expect(allTerritoryIds.has(tid)).toBe(true);
        });
      });
    });

    it('all 42 territories are assigned to exactly one continent', () => {
      const allGrouped: string[] = [];
      Object.values(continentTerritories).forEach((territories) => {
        allGrouped.push(...territories);
      });
      expect(allGrouped).toHaveLength(42);
      expect(new Set(allGrouped).size).toBe(42);
    });

    it('continent sizes match classic Risk', () => {
      expect(continentTerritories.north_america).toHaveLength(9);
      expect(continentTerritories.south_america).toHaveLength(4);
      expect(continentTerritories.europe).toHaveLength(7);
      expect(continentTerritories.africa).toHaveLength(6);
      expect(continentTerritories.asia).toHaveLength(12);
      expect(continentTerritories.australia).toHaveLength(4);
    });
  });
});
