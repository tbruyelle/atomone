# `x/multisig`

## Abstract

The `multisig` module provides a on-chain MultiSig transaction implementation,
heavily inspired by the multisig implementation of the SDK's `x/accounts`
module.

This module will serve as a cornestone of the AtomOne DAO system. The DAOs
member list will point to a `x/multisig` entry, allowing members to initiate
pending transactions to be executed and wait for other members to vote on them.

Naively, a `x/multisig` entry could be labelled as 'Steering DAO' or 'Treasury
DAO', allowing more specific messages to be submitted. Only governance proposal
can label or un-label a multisig, and label are uniques.

A multisig is created with a specific message that contains the members and a
new address is created for the purpose (see `makeAdress` in `x/accounts`)

Such multisig address will never match the treasury DAO address that was previouisly created, so the multisig would not be able to move the funds from this address, unless we move the funds to that addres thanks to an upgrade.

Messages can be executed thanks to the `baseapp.MsgServiceRouter`, like `x/gov`
when a proposal has passed. Also like `x/gov` the messages' signer must match
the multisig address, allowing modules to check that only the multisig
account has the permission to execute specific messages.


