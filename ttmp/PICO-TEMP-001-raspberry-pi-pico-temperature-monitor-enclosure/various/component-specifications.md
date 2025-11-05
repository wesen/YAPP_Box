---
Title: Component Specifications
Ticket: PICO-TEMP-001
Status: active
Topics:
    - enclosures
    - yappgenerator
    - hardware
DocType: specification
Intent: long-term
Owners:
    - manuel
RelatedFiles: []
ExternalSources: []
Summary: "Complete component specifications for Raspberry Pi Pico temperature monitor enclosure"
LastUpdated: 2025-11-05T16:02:48.420362041-05:00
---

# Component Specifications

## Complete Component Specifications for Enclosure Design

### 1. OLED Display Module (1.3" OLED MODULE)

**PCB Board Dimensions:**
- **Length**: 40mm
- **Width**: 30mm  
- **Thickness**: ~8mm (including components)

**Visible Display Area:**
- **Length**: 36mm
- **Width**: 17mm

**Mounting Hole Specifications:**
- **Hole positions**: Inset 1mm from each corner of the PCB
- **Effective mounting hole spacing**: 38mm x 28mm (center-to-center)

```
    ┌─────── 40mm ───────┐
    │ ○                 ○ │ ┐
    │   ┌─── 36mm ───┐   │ │
    │   │            │   │ │ 30mm
    │   │  17mm High │   │ │
    │   │   Display  │   │ │
    │   └────────────┘   │ │
    │ ○                 ○ │ ┘
    └─────────────────────┘
      ↑ 1mm inset holes
```

### 2. Push Button Switch

- **Threaded barrel length**: 45mm
- **Button cap diameter**: 28mm

```
    ┌──── 28mm ────┐
    │    Button    │
    │     Cap      │
    └──────┬───────┘
           │ ← Threaded
           │   Barrel
           │   45mm
           │
    ───────┴────────── Panel
```

### 3. Temperature Probe/Sensor

- **Probe diameter**: 8mm
- **Cable diameter**: 6mm
- **Length**: Variable (not critical for enclosure)

```
    ╔══════════════╗ ← 8mm probe diameter
    ║   Sensor     ║
    ╚══════════════╝
           │ ← 6mm cable
           │
           │
           └─── Variable length
```

### 4. Raspberry Pi Pico (v1) Board

**PCB Dimensions:**
- **Length**: 51mm
- **Width**: 21mm
- **Thickness**: 1mm

**Pin Header Specifications:**
- **Total pins**: 40 pins
- **Pin spacing**: 2.54mm (0.1") pitch
- **Pins per side**: 20 pins each side
- **Pin hole diameter**: 1mm
- **Header extends**: ~8.5mm from PCB surface (standard through-hole headers)

**Mounting Options:**
- **Four mounting holes**: 2.1mm diameter
- **Castellated edges**: Can be surface-mounted directly to a carrier PCB

```
    ┌─────── 51mm ───────┐
    │ ○               ○ │ ┐
    ┌┴┐ ┌┐ ┌┐ ┌┐ ┌┐ ┌┐ ┌┴┐ │
    │1│ ││ ││ ││ ││ ││ │40│ │
    ├─┤ ├┤ ├┤ ├┤ ├┤ ├┤ ├─┤ │ 21mm
    │2│ ││ ││ ││ ││ ││ │39│ │
    └┬┘ └┘ └┘ └┘ └┘ └┘ └┬┘ │
    │ ○               ○ │ ┘
    └───────────────────┘
      ↑ 2.1mm mounting holes
      ↑ 20 pins each side
```

### 5. Custom PCB with Socket Configuration

For your custom PCB that will accept the Pico via pin headers:

- **Female header sockets** needed for Pico connection
- **Socket depth**: ~8.5mm to accommodate standard male headers
- **PCB-to-PCB clearance**: ~10-12mm total height needed when Pico is plugged in

```
    Side View - Stacked Configuration:
    
    ┌───────────────────┐ ← Pico PCB (1mm thick)
    │ Raspberry Pi Pico │
    └┬┬┬┬┬┬┬┬┬┬┬┬┬┬┬┬┬┬┬┘
     ││││││││││││││││││││ ← Male headers (~8.5mm)
    ┌┴┴┴┴┴┴┴┴┴┴┴┴┴┴┴┴┴┴┴┐
    │ Custom PCB with   │ ← Your PCB with components
    │ Female Sockets    │
    └───────────────────┘
    
    Total Stack Height: ~12mm
```

### 6. Complete Assembly Layout

```
    Top View - Enclosure Layout:
    
    ┌─────────────────────────────────┐
    │  ┌─────Display─────┐            │
    │  │     36x17mm    │     ○       │ ← Button
    │  │                │   (28mm)    │   (45mm threaded)
    │  └────────────────┘             │
    │                                 │
    │    ┌─Pico Stack─┐               │
    │    │   51x21mm  │     ╔═══╗     │ ← Temp Probe
    │    │           │     ║ 8mm      │   (6mm cable)
    │    └───────────┘     ╚═══╝     │
    │                                 │
    └─────────────────────────────────┘
```

### Mounting Recommendations:

1. **Display**: Use the four corner mounting holes with M2 screws
2. **Button**: Mount through panel hole with threaded barrel
3. **Probe**: Strain relief grommet for 8mm probe + 6mm cable
4. **Custom PCB**: Design with mounting holes to secure the entire assembly
