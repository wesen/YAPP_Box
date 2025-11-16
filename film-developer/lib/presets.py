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


def list_developers(presets: list) -> list[str]:
    names = {p.get("developer", "") for p in presets}
    names.discard("")
    return sorted(names)


def list_films_for_developer(presets: list, developer: str) -> list[str]:
    names = {p.get("film", "") for p in presets if p.get("developer") == developer}
    names.discard("")
    return sorted(names)


def list_isos_for_film_dev(presets: list, developer: str, film: str) -> list[int]:
    vals = {int(p.get("iso")) for p in presets if p.get("developer") == developer and p.get("film") == film and p.get("iso") is not None}
    return sorted(vals)


def find_best_preset(presets: list, developer: str, film: str, iso: int):
    """
    Pick dilution 1+1 at ~20C if available, else any match.
    """
    candidates = [p for p in presets if p.get("developer") == developer and p.get("film") == film and int(p.get("iso", -1)) == iso]
    if not candidates:
        return None
    # Prefer 1+1 then stock; prefer 20C
    def pref(p):
        dil = (p.get("dilution") or "").lower()
        score_dil = 0 if dil == "1+1" else 1 if dil == "stock" else 2
        temp = float(p.get("temp_c") or 20.0)
        score_temp = abs(temp - 20.0)
        return (score_dil, score_temp)
    candidates.sort(key=pref)
    return candidates[0]


