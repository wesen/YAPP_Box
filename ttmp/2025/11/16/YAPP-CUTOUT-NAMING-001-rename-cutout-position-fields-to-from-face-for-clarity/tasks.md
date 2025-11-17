# Tasks

## Implementation Checklist

- [ ] Update cutouts schema.yaml with new fields (from_face_left, from_face_bottom, from_face_back)
- [ ] Implement CustomValidate() with face-type-specific validation and clear error messages
- [ ] Update Build() method to use new field names
- [ ] Run schemagen to regenerate schema_gen.go
- [ ] Add comprehensive test cases (18+ tests covering all faces and error conditions)
- [ ] Migrate all example YAML files (9 files, 60+ usages)
- [ ] Regenerate comparison STLs and verify no visual regressions
- [ ] Update pkg/docs/tutorials/yapp-dsl-reference.md with new field names
- [ ] Update face/axes documentation section
- [ ] Run full test suite to ensure no regressions

## Validation Test Cases Needed

- [ ] Front face with correct fields (from_face_left, from_face_bottom)
- [ ] Back face with correct fields
- [ ] Left face with correct fields
- [ ] Right face with correct fields
- [ ] Base face with correct fields (from_face_left, from_face_back)
- [ ] Lid face with correct fields (from_face_left, from_face_back)
- [ ] Front face with wrong field (from_face_back) → error
- [ ] Back face with wrong field (from_face_back) → error  
- [ ] Left face with wrong field (from_face_back) → error
- [ ] Right face with wrong field (from_face_back) → error
- [ ] Base face with wrong field (from_face_bottom) → error
- [ ] Lid face with wrong field (from_face_bottom) → error
- [ ] Missing from_face_left → error
- [ ] Missing from_face_bottom (side face) → error
- [ ] Missing from_face_back (base/lid) → error
