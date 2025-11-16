import json


def load_presets(path: str = "data/presets.json") -> list:
    try:
        with open(path, "r") as f:
            data = json.load(f)
        presets = data.get("presets", [])
        return presets
    except Exception:
        return []


def get_preset_by_id(presets: list, preset_id: str):
    for p in presets:
        if p.get("id") == preset_id:
            return p
    return None


