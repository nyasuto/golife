package core

import "testing"

func TestCoordConstructors(t *testing.T) {
	t.Run("NewCoord2D", func(t *testing.T) {
		c := NewCoord2D(10, 20)
		if c.X != 10 || c.Y != 20 || c.Z != 0 || c.W != 0 {
			t.Errorf("Expected (10,20,0,0), got (%d,%d,%d,%d)", c.X, c.Y, c.Z, c.W)
		}
	})

	t.Run("NewCoord3D", func(t *testing.T) {
		c := NewCoord3D(10, 20, 30)
		if c.X != 10 || c.Y != 20 || c.Z != 30 || c.W != 0 {
			t.Errorf("Expected (10,20,30,0), got (%d,%d,%d,%d)", c.X, c.Y, c.Z, c.W)
		}
	})

	t.Run("NewCoord4D", func(t *testing.T) {
		c := NewCoord4D(10, 20, 30, 40)
		if c.X != 10 || c.Y != 20 || c.Z != 30 || c.W != 40 {
			t.Errorf("Expected (10,20,30,40), got (%d,%d,%d,%d)", c.X, c.Y, c.Z, c.W)
		}
	})
}

func TestCellState(t *testing.T) {
	t.Run("Dead constant", func(t *testing.T) {
		if Dead != 0 {
			t.Errorf("Dead should be 0, got %d", Dead)
		}
	})

	t.Run("Alive constant", func(t *testing.T) {
		if Alive != 255 {
			t.Errorf("Alive should be 255, got %d", Alive)
		}
	})

	t.Run("Intermediate states", func(t *testing.T) {
		var state CellState = 128
		if state < Dead || state > Alive {
			t.Errorf("State %d should be between Dead and Alive", state)
		}
	})
}

func TestDimension(t *testing.T) {
	tests := []struct {
		name string
		dim  Dimension
		want int
	}{
		{"2D", Dim2D, 2},
		{"2.5D", Dim25D, 25},
		{"3D", Dim3D, 3},
		{"4D", Dim4D, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.dim) != tt.want {
				t.Errorf("Expected %d, got %d", tt.want, int(tt.dim))
			}
		})
	}
}

func TestNeighborhoodType(t *testing.T) {
	tests := []struct {
		name string
		nt   NeighborhoodType
		want string
	}{
		{"Moore", Moore, "Moore"},
		{"VonNeumann", VonNeumann, "VonNeumann"},
		{"Custom", Custom, "Custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.nt.String()
			if got != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, got)
			}
		})
	}
}
