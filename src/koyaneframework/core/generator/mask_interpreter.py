from koyaneframework.core.generator.maskChar import MaskChar


class MaskInterpreter:
    mask: str = None
    mask_segments: list[MaskChar] = []  # Python 3.9+

    def __init__(self, mask: str):
        self.mask = mask.strip()
        self.mask_segments = []

        segment = ""
        inside_mask = False

        for char in self.mask:
            if char == "?":
                if inside_mask:  # add built segment
                    self.mask_segments.append(MaskChar(segment))
                    segment = "?"
                else:
                    segment = "?"
                    inside_mask = True
            else:
                segment += char

        if segment.startswith("?"):  # add last segment
            self.mask_segments.append(MaskChar(segment))
