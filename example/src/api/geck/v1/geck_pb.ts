// @generated by protoc-gen-es v1.6.0 with parameter "target=ts"
// @generated from file geck/v1/geck.proto (package natsproxy.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3, protoInt64 } from "@bufbuild/protobuf";

/**
 * @generated from enum natsproxy.v1.KnownID
 */
export enum KnownID {
  /**
   * @generated from enum value: UNKNOWN = 0;
   */
  UNKNOWN = 0,

  /**
   * @generated from enum value: INTERNAL = 1;
   */
  INTERNAL = 1,

  /**
   * @generated from enum value: INDENTIFIER = 2;
   */
  INDENTIFIER = 2,

  /**
   * @generated from enum value: NAME = 3;
   */
  NAME = 3,

  /**
   * @generated from enum value: WILDCARD = 4;
   */
  WILDCARD = 4,

  /**
   * @generated from enum value: CHILD_OF = 5;
   */
  CHILD_OF = 5,

  /**
   * @generated from enum value: INSTANCE_OF = 6;
   */
  INSTANCE_OF = 6,

  /**
   * @generated from enum value: COMPONENT = 7;
   */
  COMPONENT = 7,

  /**
   * @generated from enum value: USER_DEFINED = 1000;
   */
  USER_DEFINED = 1000,
}
// Retrieve enum metadata with: proto3.getEnumType(KnownID)
proto3.util.setEnumType(KnownID, "natsproxy.v1.KnownID", [
  { no: 0, name: "UNKNOWN" },
  { no: 1, name: "INTERNAL" },
  { no: 2, name: "INDENTIFIER" },
  { no: 3, name: "NAME" },
  { no: 4, name: "WILDCARD" },
  { no: 5, name: "CHILD_OF" },
  { no: 6, name: "INSTANCE_OF" },
  { no: 7, name: "COMPONENT" },
  { no: 1000, name: "USER_DEFINED" },
]);

/**
 * @generated from message natsproxy.v1.ComponentColumnDefinition
 */
export class ComponentColumnDefinition extends Message<ComponentColumnDefinition> {
  /**
   * @generated from field: uint64 component_id = 1;
   */
  componentId = protoInt64.zero;

  /**
   * @generated from field: uint32 archetype_index = 2;
   */
  archetypeIndex = 0;

  /**
   * @generated from field: uint32 count = 3;
   */
  count = 0;

  /**
   * @generated from field: bytes data = 4;
   */
  data = new Uint8Array(0);

  constructor(data?: PartialMessage<ComponentColumnDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ComponentColumnDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "component_id", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 2, name: "archetype_index", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "count", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 4, name: "data", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ComponentColumnDefinition {
    return new ComponentColumnDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ComponentColumnDefinition {
    return new ComponentColumnDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ComponentColumnDefinition {
    return new ComponentColumnDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: ComponentColumnDefinition | PlainMessage<ComponentColumnDefinition> | undefined, b: ComponentColumnDefinition | PlainMessage<ComponentColumnDefinition> | undefined): boolean {
    return proto3.util.equals(ComponentColumnDefinition, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.ArchetypeDefinition
 */
export class ArchetypeDefinition extends Message<ArchetypeDefinition> {
  /**
   * @generated from field: uint64 hash = 1;
   */
  hash = protoInt64.zero;

  /**
   * @generated from field: uint32 depth = 2;
   */
  depth = 0;

  /**
   * @generated from field: repeated uint64 component_ids = 3;
   */
  componentIds: bigint[] = [];

  /**
   * @generated from field: repeated natsproxy.v1.ComponentColumnDefinition data_columns = 4;
   */
  dataColumns: ComponentColumnDefinition[] = [];

  /**
   * @generated from field: map<uint64, natsproxy.v1.ArchetypeDefinition.Edge> edges = 5;
   */
  edges: { [key: string]: ArchetypeDefinition_Edge } = {};

  /**
   * @generated from field: repeated uint64 entities = 6;
   */
  entities: bigint[] = [];

  constructor(data?: PartialMessage<ArchetypeDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ArchetypeDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "hash", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 2, name: "depth", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
    { no: 3, name: "component_ids", kind: "scalar", T: 4 /* ScalarType.UINT64 */, repeated: true },
    { no: 4, name: "data_columns", kind: "message", T: ComponentColumnDefinition, repeated: true },
    { no: 5, name: "edges", kind: "map", K: 4 /* ScalarType.UINT64 */, V: {kind: "message", T: ArchetypeDefinition_Edge} },
    { no: 6, name: "entities", kind: "scalar", T: 4 /* ScalarType.UINT64 */, repeated: true },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ArchetypeDefinition {
    return new ArchetypeDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ArchetypeDefinition {
    return new ArchetypeDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ArchetypeDefinition {
    return new ArchetypeDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: ArchetypeDefinition | PlainMessage<ArchetypeDefinition> | undefined, b: ArchetypeDefinition | PlainMessage<ArchetypeDefinition> | undefined): boolean {
    return proto3.util.equals(ArchetypeDefinition, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.ArchetypeDefinition.Edge
 */
export class ArchetypeDefinition_Edge extends Message<ArchetypeDefinition_Edge> {
  /**
   * @generated from field: uint64 add_id = 1;
   */
  addId = protoInt64.zero;

  /**
   * @generated from field: uint64 remove_id = 2;
   */
  removeId = protoInt64.zero;

  constructor(data?: PartialMessage<ArchetypeDefinition_Edge>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ArchetypeDefinition.Edge";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "add_id", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 2, name: "remove_id", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ArchetypeDefinition_Edge {
    return new ArchetypeDefinition_Edge().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ArchetypeDefinition_Edge {
    return new ArchetypeDefinition_Edge().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ArchetypeDefinition_Edge {
    return new ArchetypeDefinition_Edge().fromJsonString(jsonString, options);
  }

  static equals(a: ArchetypeDefinition_Edge | PlainMessage<ArchetypeDefinition_Edge> | undefined, b: ArchetypeDefinition_Edge | PlainMessage<ArchetypeDefinition_Edge> | undefined): boolean {
    return proto3.util.equals(ArchetypeDefinition_Edge, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.IDRecordDefinition
 */
export class IDRecordDefinition extends Message<IDRecordDefinition> {
  /**
   * @generated from field: natsproxy.v1.ArchetypeDefinition archetype = 1;
   */
  archetype?: ArchetypeDefinition;

  /**
   * @generated from field: uint32 row = 2;
   */
  row = 0;

  constructor(data?: PartialMessage<IDRecordDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.IDRecordDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "archetype", kind: "message", T: ArchetypeDefinition },
    { no: 2, name: "row", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): IDRecordDefinition {
    return new IDRecordDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): IDRecordDefinition {
    return new IDRecordDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): IDRecordDefinition {
    return new IDRecordDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: IDRecordDefinition | PlainMessage<IDRecordDefinition> | undefined, b: IDRecordDefinition | PlainMessage<IDRecordDefinition> | undefined): boolean {
    return proto3.util.equals(IDRecordDefinition, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.ArchetypeToRowMap
 */
export class ArchetypeToRowMap extends Message<ArchetypeToRowMap> {
  /**
   * @generated from field: map<uint64, uint32> value = 1;
   */
  value: { [key: string]: number } = {};

  constructor(data?: PartialMessage<ArchetypeToRowMap>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ArchetypeToRowMap";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "value", kind: "map", K: 4 /* ScalarType.UINT64 */, V: {kind: "scalar", T: 13 /* ScalarType.UINT32 */} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ArchetypeToRowMap {
    return new ArchetypeToRowMap().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ArchetypeToRowMap {
    return new ArchetypeToRowMap().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ArchetypeToRowMap {
    return new ArchetypeToRowMap().fromJsonString(jsonString, options);
  }

  static equals(a: ArchetypeToRowMap | PlainMessage<ArchetypeToRowMap> | undefined, b: ArchetypeToRowMap | PlainMessage<ArchetypeToRowMap> | undefined): boolean {
    return proto3.util.equals(ArchetypeToRowMap, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.ComponentToArchetype
 */
export class ComponentToArchetype extends Message<ComponentToArchetype> {
  /**
   * @generated from field: map<uint64, natsproxy.v1.ArchetypeToRowMap> value = 1;
   */
  value: { [key: string]: ArchetypeToRowMap } = {};

  constructor(data?: PartialMessage<ComponentToArchetype>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ComponentToArchetype";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "value", kind: "map", K: 4 /* ScalarType.UINT64 */, V: {kind: "message", T: ArchetypeToRowMap} },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ComponentToArchetype {
    return new ComponentToArchetype().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ComponentToArchetype {
    return new ComponentToArchetype().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ComponentToArchetype {
    return new ComponentToArchetype().fromJsonString(jsonString, options);
  }

  static equals(a: ComponentToArchetype | PlainMessage<ComponentToArchetype> | undefined, b: ComponentToArchetype | PlainMessage<ComponentToArchetype> | undefined): boolean {
    return proto3.util.equals(ComponentToArchetype, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.ComponentMetadataDefinition
 */
export class ComponentMetadataDefinition extends Message<ComponentMetadataDefinition> {
  /**
   * @generated from field: uint64 id = 1;
   */
  id = protoInt64.zero;

  /**
   * @generated from field: string name = 2;
   */
  name = "";

  /**
   * @generated from field: bytes resetExample = 3;
   */
  resetExample = new Uint8Array(0);

  /**
   * @generated from field: uint32 element_size = 4;
   */
  elementSize = 0;

  constructor(data?: PartialMessage<ComponentMetadataDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.ComponentMetadataDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "id", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 2, name: "name", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "resetExample", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
    { no: 4, name: "element_size", kind: "scalar", T: 13 /* ScalarType.UINT32 */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): ComponentMetadataDefinition {
    return new ComponentMetadataDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): ComponentMetadataDefinition {
    return new ComponentMetadataDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): ComponentMetadataDefinition {
    return new ComponentMetadataDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: ComponentMetadataDefinition | PlainMessage<ComponentMetadataDefinition> | undefined, b: ComponentMetadataDefinition | PlainMessage<ComponentMetadataDefinition> | undefined): boolean {
    return proto3.util.equals(ComponentMetadataDefinition, a, b);
  }
}

/**
 * @generated from message natsproxy.v1.WorldDefinition
 */
export class WorldDefinition extends Message<WorldDefinition> {
  /**
   * @generated from field: repeated uint64 available_id = 1;
   */
  availableId: bigint[] = [];

  /**
   * @generated from field: uint64 next_id = 2;
   */
  nextId = protoInt64.zero;

  /**
   * @generated from field: map<uint64, natsproxy.v1.ComponentMetadataDefinition> component_metadata = 3;
   */
  componentMetadata: { [key: string]: ComponentMetadataDefinition } = {};

  /**
   * @generated from field: map<uint64, natsproxy.v1.ArchetypeDefinition> archetypes = 4;
   */
  archetypes: { [key: string]: ArchetypeDefinition } = {};

  /**
   * @generated from field: natsproxy.v1.ComponentToArchetype archetype_component_comlumn_indicies = 5;
   */
  archetypeComponentComlumnIndicies?: ComponentToArchetype;

  constructor(data?: PartialMessage<WorldDefinition>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "natsproxy.v1.WorldDefinition";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "available_id", kind: "scalar", T: 4 /* ScalarType.UINT64 */, repeated: true },
    { no: 2, name: "next_id", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 3, name: "component_metadata", kind: "map", K: 4 /* ScalarType.UINT64 */, V: {kind: "message", T: ComponentMetadataDefinition} },
    { no: 4, name: "archetypes", kind: "map", K: 4 /* ScalarType.UINT64 */, V: {kind: "message", T: ArchetypeDefinition} },
    { no: 5, name: "archetype_component_comlumn_indicies", kind: "message", T: ComponentToArchetype },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): WorldDefinition {
    return new WorldDefinition().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): WorldDefinition {
    return new WorldDefinition().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): WorldDefinition {
    return new WorldDefinition().fromJsonString(jsonString, options);
  }

  static equals(a: WorldDefinition | PlainMessage<WorldDefinition> | undefined, b: WorldDefinition | PlainMessage<WorldDefinition> | undefined): boolean {
    return proto3.util.equals(WorldDefinition, a, b);
  }
}

