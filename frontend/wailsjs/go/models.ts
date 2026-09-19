export namespace api {
	
	export class DepositView {
	    id: number;
	    type: string;
	    total: number;
	    remaining: number;
	
	    static createFrom(source: any = {}) {
	        return new DepositView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.total = source["total"];
	        this.remaining = source["remaining"];
	    }
	}
	export class DroneView {
	    id: number;
	    faction_id: number;
	    name: string;
	    type: string;
	    level: number;
	    hp: number;
	    attack: number;
	    activity: string;
	    target: number;
	
	    static createFrom(source: any = {}) {
	        return new DroneView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.faction_id = source["faction_id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.level = source["level"];
	        this.hp = source["hp"];
	        this.attack = source["attack"];
	        this.activity = source["activity"];
	        this.target = source["target"];
	    }
	}
	export class FactionView {
	    id: number;
	    name: string;
	    territories: number[];
	    drones: number[];
	    resources: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new FactionView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.territories = source["territories"];
	        this.drones = source["drones"];
	        this.resources = source["resources"];
	    }
	}
	export class TerritoryView {
	    id: number;
	    owner: number;
	    factions: number[];
	    deposit_types: string[];
	    deposits: number[];
	
	    static createFrom(source: any = {}) {
	        return new TerritoryView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.owner = source["owner"];
	        this.factions = source["factions"];
	        this.deposit_types = source["deposit_types"];
	        this.deposits = source["deposits"];
	    }
	}

}

