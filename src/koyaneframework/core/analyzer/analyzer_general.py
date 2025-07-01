from pathlib import Path
import datetime
import chardet
import hashlib

class GeneralAnalyzer:
    file_path: Path = None

    def __init__(self, file_path: Path):
        self.file_path = file_path




    def name(self):
        return self.file_path.name

    def path(self):
        return self.file_path.resolve()

    def size(self):
        return self.file_path.stat().st_size

    def last_changes(self):
        return datetime.datetime.fromtimestamp(self.file_path.stat().st_mtime)

    def last_access(self):
        return datetime.datetime.fromtimestamp(self.file_path.stat().st_atime)

    def created(self):
        return datetime.datetime.fromtimestamp(self.file_path.stat().st_ctime)

    def extension(self):
        return self.file_path.suffix

    def encoding(self, sample_size: int = 100_000) -> str:
        with self.file_path.open("rb") as f:
            sample = f.read(sample_size)  #  100.000 Bytes
        result = chardet.detect(sample)
        return result.get("encoding", "Unknown")

    def hash(self):
        hash_md5 = hashlib.md5()
        with self.file_path.open("rb") as f:
            for chunk in iter(lambda: f.read(4096), b""):
                hash_md5.update(chunk)
        return hash_md5.hexdigest()


