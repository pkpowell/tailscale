class AppData {
    Profile: Profile;
    DeviceName: string;
    NodeKey: string;
    StableID: string;
    Created: string;
    ServerURL: string;
    Version: string;
    Arch: string;
    OS: string;
    OSVersion: string;
    HostName: string;
    Name: string;
    Owner: string;
    Services: Service[];
    Health: string[];
    IPs: object[];
    IPv4: string;
    IPv6: string;
    PeerAPIPort: number;
    Now: Date;
    TX: number;
    RX: number;
    Peers: Peer[];
}

class Factory {
    create<T>(type: (new () => T)): T {
        return new type();
    }
}

let factory = new Factory()
let Appl = factory.create(AppData)

// interface PeerHash {

// }
interface Service {
    Description: string;
    Port: number;
    Proto: string;
}

interface Profile {
    ID: string;
    LoginName: string;
    DisplayName: string;
    ProfilePicURL: string;
    Roles: object[];
}

interface Peer {
    HostName: string;
    ID: string;
    UserID: string;
    NodeKey: string;
    OS: string;
    Created: Date;
    DNSName: string;
    IPs: string[];
    TailscaleIPs: string[];
    IPv4: string;
    IPv6: string;
    IPv4Num: string;
    IPv6Num: string;
    PeerAPIPort: number;
    RX: string;
    TX: string;
    RxBytes: number;
    TxBytes: number;
    RXb: number;
    TXb: number;
    Relay: string;
    RelayActive: boolean;
    Connection: string;
    ActAgo: string;
    LastSeen: Date;
    Unseen: boolean;
    CreatedDate: Date;
    CreatedTime: Date;
}

interface Base {
    factor: number; 
    suffix: string;
}

interface BytePrefix {
    short: string;
    full: string;
}
interface SSEMessage {
    length: number;
    type: string;
    data: any;
    payload: any;
}

type BPrefix = Record<number, BytePrefix>

export type {
    Peer,
    AppData,
    Appl,
    Base,
    BPrefix,
    SSEMessage
}