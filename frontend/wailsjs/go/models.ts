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

