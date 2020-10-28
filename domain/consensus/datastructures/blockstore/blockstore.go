package blockstore

import (
	"github.com/kaspanet/kaspad/domain/consensus/model"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/dbkeys"
)

var bucket = dbkeys.MakeBucket([]byte("blocks"))

// blockStore represents a store of blocks
type blockStore struct {
	staging map[externalapi.DomainHash]*externalapi.DomainBlock
}

// New instantiates a new BlockStore
func New() model.BlockStore {
	return &blockStore{
		staging: make(map[externalapi.DomainHash]*externalapi.DomainBlock),
	}
}

// Stage stages the given block for the given blockHash
func (bms *blockStore) Stage(blockHash *externalapi.DomainHash, block *externalapi.DomainBlock) {
	bms.staging[*blockHash] = block
}

func (bms *blockStore) IsStaged() bool {
	return len(bms.staging) != 0
}

func (bms *blockStore) Discard() {
	bms.staging = make(map[externalapi.DomainHash]*externalapi.DomainBlock)
}

func (bms *blockStore) Commit(dbTx model.DBTransaction) error {
	for hash, block := range bms.staging {
		err := dbTx.Put(bms.hashAsKey(&hash), bms.serializeBlock(block))
		if err != nil {
			return err
		}
	}

	bms.Discard()
	return nil
}

// Block gets the block associated with the given blockHash
func (bms *blockStore) Block(dbContext model.DBReader, blockHash *externalapi.DomainHash) (*externalapi.DomainBlock, error) {
	if block, ok := bms.staging[*blockHash]; ok {
		return block, nil
	}

	blockBytes, err := dbContext.Get(bms.hashAsKey(blockHash))
	if err != nil {
		return nil, err
	}

	return bms.deserializeBlock(blockBytes)
}

// HasBlock returns whether a block with a given hash exists in the store.
func (bms *blockStore) HasBlock(dbContext model.DBReader, blockHash *externalapi.DomainHash) (bool, error) {
	if _, ok := bms.staging[*blockHash]; ok {
		return true, nil
	}

	exists, err := dbContext.Has(bms.hashAsKey(blockHash))
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Blocks gets the blocks associated with the given blockHashes
func (bms *blockStore) Blocks(dbContext model.DBReader, blockHashes []*externalapi.DomainHash) ([]*externalapi.DomainBlock, error) {
	blocks := make([]*externalapi.DomainBlock, len(blockHashes))
	for i, hash := range blockHashes {
		var err error
		blocks[i], err = bms.Block(dbContext, hash)
		if err != nil {
			return nil, err
		}
	}
	return blocks, nil
}

// Delete deletes the block associated with the given blockHash
func (bms *blockStore) Delete(dbTx model.DBTransaction, blockHash *externalapi.DomainHash) error {
	if _, ok := bms.staging[*blockHash]; ok {
		delete(bms.staging, *blockHash)
		return nil
	}
	return dbTx.Delete(bms.hashAsKey(blockHash))
}

func (bms *blockStore) serializeBlock(block *externalapi.DomainBlock) []byte {
	panic("implement me")
}

func (bms *blockStore) deserializeBlock(blockBytes []byte) (*externalapi.DomainBlock, error) {
	panic("implement me")
}

func (bms *blockStore) hashAsKey(hash *externalapi.DomainHash) model.DBKey {
	return bucket.Key(hash[:])
}
