/** Represents a configuration that only this module can refer to. */
//dokimi:builder
class Config {
  /** Contains the name of the configuration. */
  name: string = "";
}

/** Returns the name of a new configuration. */
export function configName(): string {
  return new Config().name;
}
