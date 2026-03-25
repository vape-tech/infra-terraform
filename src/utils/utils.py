import os
import json
import logging
from typing import Dict, Any, Optional

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def load_config(config_file: str) -> Dict[str, Any]:
    if not os.path.exists(config_file):
        raise FileNotFoundError(f"Config file '{config_file}' not found.")
    
    with open(config_file, 'r') as file:
        return json.load(file)

def save_config(config: Dict[str, Any], config_file: str) -> None:
    with open(config_file, 'w') as file:
        json.dump(config, file, indent=4)

def validate_config(config: Dict[str, Any], required_keys: list) -> bool:
    missing_keys = [key for key in required_keys if key not in config]
    if missing_keys:
        logging.error(f"Missing required keys in config: {missing_keys}")
        return False
    return True

def create_directory_if_not_exists(directory: str) -> None:
    if not os.path.exists(directory):
        os.makedirs(directory)
        logging.info(f"Created directory: {directory}")

def get_env_variable(key: str, default: Optional[str] = None) -> str:
    value = os.getenv(key, default)
    if value is None:
        raise EnvironmentError(f"Environment variable '{key}' is not set.")
    return value